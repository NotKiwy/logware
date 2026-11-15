package subs

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Sub struct {
	Plan   string
	Until  time.Time
	Charge string
}

func (s *Sub) Live() bool {
	return s != nil && time.Now().Before(s.Until)
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Get(__ctx context.Context, __user int64) (*Sub, error) {
	_s := &Sub{}
	_err := r.__db.QueryRow(__ctx, `select plan, until, charge from subs where user_id = $1`, __user).Scan(&_s.Plan, &_s.Until, &_s.Charge)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _s, _err
}

func (r *Repo) Pay(__ctx context.Context, __user int64, __plan string, __amount int, __cur string, __until time.Time, __charge string) (bool, error) {
	_tx, _err := r.__db.Begin(__ctx)
	if _err != nil {
		return false, _err
	}
	defer _tx.Rollback(__ctx)
	_tag, _err := _tx.Exec(__ctx, `insert into pays (charge, user_id, plan, amount, until, cur) values ($1, $2, $3, $4, $5, $6) on conflict do nothing`,
		__charge, __user, __plan, __amount, __until, __cur)
	if _err != nil || _tag.RowsAffected() == 0 {
		return false, _err
	}
	if _, _err = _tx.Exec(__ctx, `
		insert into subs (user_id, plan, until, charge) values ($1, $2, $3, $4)
		on conflict (user_id) do update set plan = excluded.plan, until = excluded.until, charge = excluded.charge, warned = null, updated = now()`,
		__user, __plan, __until, __charge); _err != nil {
		return false, _err
	}
	return true, _tx.Commit(__ctx)
}

func (r *Repo) Take(__ctx context.Context, __user int64) (string, error) {
	var _c string
	_err := r.__db.QueryRow(__ctx, `delete from subs where user_id = $1 returning charge`, __user).Scan(&_c)
	if errors.Is(_err, pgx.ErrNoRows) {
		return "", nil
	}
	return _c, _err
}

type Due struct {
	User  int64
	Plan  string
	Until time.Time
}

func (r *Repo) Due(__ctx context.Context, __in time.Duration) ([]Due, error) {
	_rows, _err := r.__db.Query(__ctx, `
		update subs set warned = now()
		where warned is null and until > now() and until < now() + $1::interval and (charge like 'cb:%' or charge like 'xr:%')
		returning user_id, plan, until`, __in.String())
	if _err != nil {
		return nil, _err
	}
	defer _rows.Close()
	var _out []Due
	for _rows.Next() {
		var _d Due
		if _err = _rows.Scan(&_d.User, &_d.Plan, &_d.Until); _err != nil {
			return nil, _err
		}
		_out = append(_out, _d)
	}
	return _out, _rows.Err()
}
