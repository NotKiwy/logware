package ign

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Has(__ctx context.Context, __user, __chat int64) (bool, error) {
	var _ok bool
	_err := r.__db.QueryRow(__ctx, `select exists (select 1 from ignores where user_id = $1 and chat_id = $2)`, __user, __chat).Scan(&_ok)
	return _ok, _err
}

func (r *Repo) List(__ctx context.Context, __user int64) (map[int64]bool, error) {
	_rs, _err := r.__db.Query(__ctx, `select chat_id from ignores where user_id = $1`, __user)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	_out := map[int64]bool{}
	for _rs.Next() {
		var _c int64
		if _err = _rs.Scan(&_c); _err != nil {
			return nil, _err
		}
		_out[_c] = true
	}
	return _out, _rs.Err()
}

func (r *Repo) Set(__ctx context.Context, __user, __chat int64, __on bool) error {
	if __on {
		_, _err := r.__db.Exec(__ctx, `insert into ignores (user_id, chat_id) values ($1, $2) on conflict do nothing`, __user, __chat)
		return _err
	}
	_, _err := r.__db.Exec(__ctx, `delete from ignores where user_id = $1 and chat_id = $2`, __user, __chat)
	return _err
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from ignores where user_id = $1`, __user)
	return _err
}
