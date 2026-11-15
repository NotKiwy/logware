package promos

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNone = errors.New("no code")
	ErrUsed = errors.New("used")
	ErrOver = errors.New("over")
)

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func Norm(__s string) string {
	return strings.ToUpper(strings.TrimSpace(__s))
}

type Merge func(string, time.Time, string, time.Duration, time.Time) (string, time.Time)

func (r *Repo) Use(__ctx context.Context, __code string, __user int64, __merge Merge) (string, time.Time, error) {
	_tx, _err := r.__db.Begin(__ctx)
	if _err != nil {
		return "", time.Time{}, _err
	}
	defer _tx.Rollback(__ctx)
	if _, _err = _tx.Exec(__ctx, `select pg_advisory_xact_lock(7, $1::int)`, __user%2147483647); _err != nil {
		return "", time.Time{}, _err
	}
	var _plan string
	var _days, _cap, _used int
	var _ends *time.Time
	var _off bool
	_err = _tx.QueryRow(__ctx, `select plan, days, cap, used, ends, off from promos where code = $1 for update`, Norm(__code)).Scan(&_plan, &_days, &_cap, &_used, &_ends, &_off)
	if errors.Is(_err, pgx.ErrNoRows) || (_err == nil && _off) {
		return "", time.Time{}, ErrNone
	}
	if _err != nil {
		return "", time.Time{}, _err
	}
	if (_cap > 0 && _used >= _cap) || (_ends != nil && time.Now().After(*_ends)) {
		return "", time.Time{}, ErrOver
	}
	_t, _err := _tx.Exec(__ctx, `insert into promo_uses (code, user_id) values ($1, $2) on conflict do nothing`, Norm(__code), __user)
	if _err != nil {
		return "", time.Time{}, _err
	}
	if _t.RowsAffected() == 0 {
		return "", time.Time{}, ErrUsed
	}
	if _, _err = _tx.Exec(__ctx, `update promos set used = used + 1 where code = $1`, Norm(__code)); _err != nil {
		return "", time.Time{}, _err
	}
	var _cur string
	var _until time.Time
	_err = _tx.QueryRow(__ctx, `select plan, until from subs where user_id = $1 and until > now()`, __user).Scan(&_cur, &_until)
	if _err != nil && !errors.Is(_err, pgx.ErrNoRows) {
		return "", time.Time{}, _err
	}
	if _err != nil {
		_cur, _until = "", time.Time{}
	}
	_plan, _until = __merge(_cur, _until, _plan, time.Duration(_days)*24*time.Hour, time.Now())
	if _, _err = _tx.Exec(__ctx, `
		insert into subs (user_id, plan, until, charge) values ($1, $2, $3, $4)
		on conflict (user_id) do update set plan = excluded.plan, until = excluded.until, warned = null, updated = now(),
			charge = case when subs.until <= now() or subs.charge like 'admin%' or subs.charge like 'promo:%' then excluded.charge else subs.charge end`,
		__user, _plan, _until, "promo:"+Norm(__code)); _err != nil {
		return "", time.Time{}, _err
	}
	return _plan, _until, _tx.Commit(__ctx)
}
