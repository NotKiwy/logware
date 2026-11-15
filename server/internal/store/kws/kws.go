package kws

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const Max = 20

type Kw struct {
	Hash []byte `json:"hash"`
	Box  []byte `json:"box"`
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) List(__ctx context.Context, __user int64) ([]Kw, error) {
	_rs, _err := r.__db.Query(__ctx, `select hash, box from kws where user_id = $1 order by created`, __user)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	_out := []Kw{}
	for _rs.Next() {
		var _k Kw
		if _err = _rs.Scan(&_k.Hash, &_k.Box); _err != nil {
			return nil, _err
		}
		_out = append(_out, _k)
	}
	return _out, _rs.Err()
}

func (r *Repo) Add(__ctx context.Context, __user int64, __k Kw) (bool, error) {
	_tag, _err := r.__db.Exec(__ctx, `
		insert into kws (user_id, hash, box)
		select $1, $2, $3 where (select count(*) from kws where user_id = $1) < $4
		on conflict do nothing`, __user, __k.Hash, __k.Box, Max)
	return _tag.RowsAffected() > 0, _err
}

func (r *Repo) Del(__ctx context.Context, __user int64, __hash []byte) error {
	_, _err := r.__db.Exec(__ctx, `delete from kws where user_id = $1 and hash = $2`, __user, __hash)
	return _err
}

func (r *Repo) Hit(__ctx context.Context, __user int64, __hashes [][]byte) (bool, error) {
	if len(__hashes) == 0 {
		return false, nil
	}
	var _ok bool
	_err := r.__db.QueryRow(__ctx, `select exists (select 1 from kws where user_id = $1 and hash = any($2))`, __user, __hashes).Scan(&_ok)
	return _ok, _err
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from kws where user_id = $1`, __user)
	return _err
}
