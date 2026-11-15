package devs

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaken = errors.New("nick taken")

type Dev struct {
	User int64
	Nick string
	Bio  string
}

type Req struct {
	ID      int64
	User    int64
	Name    string
	About   string
	Ver     string
	Note    string
	Repo    string
	Plug    int64
	Kind    string
	Pic     bool
	Status  string
	Reason  string
	Msg     int64
	Created time.Time
	Done    time.Time
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Get(__ctx context.Context, __user int64) (*Dev, error) {
	_d := &Dev{User: __user}
	_err := r.__db.QueryRow(__ctx, `select nick, bio from devs where user_id = $1`, __user).Scan(&_d.Nick, &_d.Bio)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _d, _err
}

func (r *Repo) Nicks(__ctx context.Context, __ids []int64) (map[int64]string, error) {
	_out := map[int64]string{}
	if len(__ids) == 0 {
		return _out, nil
	}
	_rows, _err := r.__db.Query(__ctx, `select user_id, nick from devs where user_id = any($1)`, __ids)
	if _err != nil {
		return nil, _err
	}
	defer _rows.Close()
	for _rows.Next() {
		var _id int64
		var _n string
		if _err = _rows.Scan(&_id, &_n); _err != nil {
			return nil, _err
		}
		_out[_id] = _n
	}
	return _out, _rows.Err()
}

func (r *Repo) Set(__ctx context.Context, __d Dev) error {
	_, _err := r.__db.Exec(__ctx, `
		insert into devs (user_id, nick, bio) values ($1, $2, $3)
		on conflict (user_id) do update set nick = excluded.nick, bio = excluded.bio`, __d.User, __d.Nick, __d.Bio)
	var _pe *pgconn.PgError
	if errors.As(_err, &_pe) && _pe.Code == "23505" {
		return ErrTaken
	}
	return _err
}

const __cols__ = `id, user_id, name, about, ver, note, repo, coalesce(plug_id, 0), kind, pic is not null, status, reason, msg, created, coalesce(done, created)`

func scan(__r pgx.Row, __q *Req) error {
	return __r.Scan(&__q.ID, &__q.User, &__q.Name, &__q.About, &__q.Ver, &__q.Note, &__q.Repo, &__q.Plug, &__q.Kind, &__q.Pic, &__q.Status, &__q.Reason, &__q.Msg, &__q.Created, &__q.Done)
}

func (r *Repo) Req(__ctx context.Context, __id int64) (*Req, error) {
	_q := &Req{}
	_err := scan(r.__db.QueryRow(__ctx, `select `+__cols__+` from reqs where id = $1`, __id), _q)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _q, _err
}

func (r *Repo) ByMsg(__ctx context.Context, __msg int64) (*Req, error) {
	_q := &Req{}
	_err := scan(r.__db.QueryRow(__ctx, `select `+__cols__+` from reqs where msg = $1 and msg <> 0`, __msg), _q)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _q, _err
}

func (r *Repo) reqs(__ctx context.Context, __sql string, __args ...any) ([]Req, error) {
	_rows, _err := r.__db.Query(__ctx, __sql, __args...)
	if _err != nil {
		return nil, _err
	}
	defer _rows.Close()
	_out := []Req{}
	for _rows.Next() {
		var _q Req
		if _err = scan(_rows, &_q); _err != nil {
			return nil, _err
		}
		_out = append(_out, _q)
	}
	return _out, _rows.Err()
}

func (r *Repo) Waiting(__ctx context.Context, __user int64) (int, error) {
	var _n int
	return _n, r.__db.QueryRow(__ctx, `select count(*) from reqs where user_id = $1 and status = 'wait'`, __user).Scan(&_n)
}

func (r *Repo) Msg(__ctx context.Context, __id, __msg int64) error {
	_, _err := r.__db.Exec(__ctx, `update reqs set msg = $2 where id = $1`, __id, __msg)
	return _err
}

func (r *Repo) Mark(__ctx context.Context, __id int64, __status, __reason string) error {
	_, _err := r.__db.Exec(__ctx, `update reqs set status = $2, reason = $3, updated = now(), done = now() where id = $1`, __id, __status, __reason)
	return _err
}
