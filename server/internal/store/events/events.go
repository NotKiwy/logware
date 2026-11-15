package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Ev struct {
	ID      string
	User    int64
	Kind    string
	Refs    []int64
	Created time.Time
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Add(__ctx context.Context, __user int64, __kind string, __refs []int64) (string, error) {
	_b := make([]byte, 12)
	if _, _err := rand.Read(_b); _err != nil {
		return "", _err
	}
	_id := hex.EncodeToString(_b)
	_, _err := r.__db.Exec(__ctx, `insert into events (id, user_id, kind, refs) values ($1, $2, $3, $4)`, _id, __user, __kind, __refs)
	return _id, _err
}

func (r *Repo) Get(__ctx context.Context, __id string, __user int64) (*Ev, error) {
	_e := &Ev{}
	_err := r.__db.QueryRow(__ctx, `select id, user_id, kind, refs, created from events where id = $1 and user_id = $2`, __id, __user).
		Scan(&_e.ID, &_e.User, &_e.Kind, &_e.Refs, &_e.Created)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _e, _err
}

type Row struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	At   int64  `json:"at"`
	N    int    `json:"n"`
	Body []byte `json:"body"`
	Type string `json:"type"`
}

func (r *Repo) List(__ctx context.Context, __user int64, __off, __lim int) ([]Row, error) {
	_rs, _err := r.__db.Query(__ctx, `
		select e.id, e.kind, e.created, cardinality(e.refs), m.body, m.kind
		from events e join lateral (select body, kind from msgs where id = e.refs[cardinality(e.refs)] and user_id = e.user_id) m on true
		where e.user_id = $1 order by e.created desc offset $2 limit $3`, __user, __off, __lim)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	_out := []Row{}
	for _rs.Next() {
		var _r Row
		var _at time.Time
		if _err = _rs.Scan(&_r.ID, &_r.Kind, &_at, &_r.N, &_r.Body, &_r.Type); _err != nil {
			return nil, _err
		}
		_r.At = _at.Unix()
		_out = append(_out, _r)
	}
	return _out, _rs.Err()
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from events where user_id = $1`, __user)
	return _err
}

func (r *Repo) Clean(__ctx context.Context) (int64, error) {
	_tag, _err := r.__db.Exec(__ctx, `delete from events e where not exists (select 1 from msgs m where m.id = any(e.refs))`)
	return _tag.RowsAffected(), _err
}

func (r *Repo) Count(__ctx context.Context, __user int64) (int, int, error) {
	var _all, _day int
	_err := r.__db.QueryRow(__ctx, `select count(*), count(*) filter (where created > now() - interval '24 hours') from events where user_id = $1`, __user).Scan(&_all, &_day)
	return _all, _day, _err
}
