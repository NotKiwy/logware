package users

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID     int64
	Mode   string
	Pub    []byte
	Vault  []byte
	Prefs  Prefs
	Banned bool
}

func (u *User) Ready() bool {
	return u != nil && len(u.Pub) == 32
}

func (u *User) Live() bool {
	return u.Ready() && !u.Banned
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Get(__ctx context.Context, __id int64) (*User, error) {
	_u := &User{ID: __id}
	var _raw []byte
	_err := r.__db.QueryRow(__ctx, `select mode, pub, vault, prefs, banned from users where id = $1`, __id).
		Scan(&_u.Mode, &_u.Pub, &_u.Vault, &_raw, &_u.Banned)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	if _err != nil {
		return nil, _err
	}
	_u.Prefs = Parse(_raw)
	return _u, nil
}

func (r *Repo) Reg(__ctx context.Context, __id int64, __mode string, __pub, __vault []byte) error {
	_, _err := r.__db.Exec(__ctx, `
		insert into users (id, mode, pub, vault) values ($1, $2, $3, $4)
		on conflict (id) do update set mode = excluded.mode, pub = excluded.pub, vault = excluded.vault`,
		__id, __mode, __pub, __vault)
	return _err
}

func (r *Repo) SetPrefs(__ctx context.Context, __id int64, __p Prefs) error {
	_raw, _err := json.Marshal(__p.Fix())
	if _err != nil {
		return _err
	}
	_, _err = r.__db.Exec(__ctx, `update users set prefs = $2 where id = $1`, __id, _raw)
	return _err
}

func (r *Repo) Drop(__ctx context.Context, __id int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from users where id = $1`, __id)
	return _err
}

func (r *Repo) Touch(__ctx context.Context, __id int64, __name, __nick string, __prem bool) (bool, error) {
	var _b bool
	_err := r.__db.QueryRow(__ctx, `
		with u as (
			update users set seen = now(), name = $2, nick = $3, prem = $4
			where id = $1 and (seen is null or seen < now() - interval '1 minute' or name <> $2 or nick <> $3 or prem <> $4)
			returning banned
		)
		select coalesce((select banned from u), (select banned from users where id = $1), false)`, __id, __name, __nick, __prem).Scan(&_b)
	return _b, _err
}
