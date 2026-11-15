package mods

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	ID   string
	On   bool
	Cfg  map[string]string
	Last *time.Time
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) List(__ctx context.Context, __user int64) (map[string]State, error) {
	_rs, _err := r.__db.Query(__ctx, `select id, active, cfg, last from mods where user_id = $1`, __user)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	_out := map[string]State{}
	for _rs.Next() {
		var _s State
		var _raw []byte
		if _err = _rs.Scan(&_s.ID, &_s.On, &_raw, &_s.Last); _err != nil {
			return nil, _err
		}
		_ = json.Unmarshal(_raw, &_s.Cfg)
		_out[_s.ID] = _s
	}
	return _out, _rs.Err()
}

func (r *Repo) Set(__ctx context.Context, __user int64, __id string, __on bool, __cfg map[string]string) error {
	_raw, _err := json.Marshal(__cfg)
	if _err != nil {
		return _err
	}
	_, _err = r.__db.Exec(__ctx, `
		insert into mods (user_id, id, active, cfg) values ($1, $2, $3, $4)
		on conflict (user_id, id) do update set active = excluded.active, cfg = excluded.cfg`,
		__user, __id, __on, _raw)
	return _err
}

type Row struct {
	User int64
	Cfg  map[string]string
}

func (r *Repo) Active(__ctx context.Context, __id string) ([]Row, error) {
	_rs, _err := r.__db.Query(__ctx, `select user_id, cfg from mods where id = $1 and active`, __id)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	var _out []Row
	for _rs.Next() {
		var _r Row
		var _raw []byte
		if _err = _rs.Scan(&_r.User, &_raw); _err != nil {
			return nil, _err
		}
		_ = json.Unmarshal(_raw, &_r.Cfg)
		_out = append(_out, _r)
	}
	return _out, _rs.Err()
}

func (r *Repo) Take(__ctx context.Context, __user int64, __id string, __win time.Duration) (bool, time.Time, error) {
	var _at time.Time
	_err := r.__db.QueryRow(__ctx, `
		update mods set last = now()
		where user_id = $1 and id = $2 and (last is null or last <= now() - make_interval(secs => $3))
		returning last`, __user, __id, __win.Seconds()).Scan(&_at)
	if _err == nil {
		return true, _at, nil
	}
	if !errors.Is(_err, pgx.ErrNoRows) {
		return false, _at, _err
	}
	_err = r.__db.QueryRow(__ctx, `select coalesce(last, now()) from mods where user_id = $1 and id = $2`, __user, __id).Scan(&_at)
	return false, _at, _err
}

func (r *Repo) Free(__ctx context.Context, __user int64, __id string) error {
	_, _err := r.__db.Exec(__ctx, `update mods set last = null where user_id = $1 and id = $2`, __user, __id)
	return _err
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from mods where user_id = $1`, __user)
	return _err
}

func (r *Repo) Users(__ctx context.Context, __ids []string) (map[string]int, error) {
	_out := map[string]int{}
	if len(__ids) == 0 {
		return _out, nil
	}
	_rows, _err := r.__db.Query(__ctx, `select id, count(*) from mods where active and id = any($1) group by id`, __ids)
	if _err != nil {
		return nil, _err
	}
	defer _rows.Close()
	for _rows.Next() {
		var _id string
		var _n int
		if _err = _rows.Scan(&_id, &_n); _err != nil {
			return nil, _err
		}
		_out[_id] = _n
	}
	return _out, _rows.Err()
}
