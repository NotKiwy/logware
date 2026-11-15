package devs

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrSlug = errors.New("slug taken")

type Plug struct {
	ID      int64
	User    int64
	Slug    string
	Name    string
	About   string
	Repo    string
	Ver     string
	Pv      int
	Live    bool
	Created time.Time
	Last    string
	Why     string
	Kind    string
}

const __pcols__ = `p.id, p.user_id, p.slug, p.name, p.about, p.repo, p.ver, p.pv, p.live, p.created,
	coalesce(q.status, ''), coalesce(q.reason, ''), coalesce(q.kind, '')`

const __plast__ = ` left join lateral (select status, reason, kind from reqs where plug_id = p.id order by created desc limit 1) q on true`

func pscan(__r pgx.Row, __p *Plug) error {
	return __r.Scan(&__p.ID, &__p.User, &__p.Slug, &__p.Name, &__p.About, &__p.Repo, &__p.Ver, &__p.Pv, &__p.Live, &__p.Created, &__p.Last, &__p.Why, &__p.Kind)
}

func (r *Repo) Plugs(__ctx context.Context, __user int64) ([]Plug, error) {
	return r.plugs(__ctx, `select `+__pcols__+` from plugs p`+__plast__+` where p.user_id = $1 order by p.created desc`, __user)
}

func (r *Repo) Live(__ctx context.Context) ([]Plug, error) {
	return r.plugs(__ctx, `select `+__pcols__+` from plugs p`+__plast__+` where p.live`)
}

func (r *Repo) plugs(__ctx context.Context, __sql string, __args ...any) ([]Plug, error) {
	_rows, _err := r.__db.Query(__ctx, __sql, __args...)
	if _err != nil {
		return nil, _err
	}
	defer _rows.Close()
	_out := []Plug{}
	for _rows.Next() {
		var _p Plug
		if _err = pscan(_rows, &_p); _err != nil {
			return nil, _err
		}
		_out = append(_out, _p)
	}
	return _out, _rows.Err()
}

func (r *Repo) Plug(__ctx context.Context, __id int64) (*Plug, error) {
	_p := &Plug{}
	_err := pscan(r.__db.QueryRow(__ctx, `select `+__pcols__+` from plugs p`+__plast__+` where p.id = $1`, __id), _p)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _p, _err
}

func (r *Repo) Pic(__ctx context.Context, __id int64) ([]byte, error) {
	var _b []byte
	_err := r.__db.QueryRow(__ctx, `select pic from plugs where id = $1 and live and pic is not null`, __id).Scan(&_b)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _b, _err
}

func (r *Repo) ReqPic(__ctx context.Context, __id int64) ([]byte, int64, error) {
	var _b []byte
	var _u int64
	_err := r.__db.QueryRow(__ctx, `select pic, user_id from reqs where id = $1`, __id).Scan(&_b, &_u)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, 0, nil
	}
	return _b, _u, _err
}

func (r *Repo) Create(__ctx context.Context, __p Plug, __q Req, __pic []byte) (int64, int64, error) {
	_tx, _err := r.__db.Begin(__ctx)
	if _err != nil {
		return 0, 0, _err
	}
	defer _tx.Rollback(__ctx)
	var _pid, _qid int64
	_err = _tx.QueryRow(__ctx, `insert into plugs (user_id, slug, name, about, repo, ver) values ($1, $2, $3, $4, $5, $6) returning id`,
		__p.User, __p.Slug, __p.Name, __p.About, __p.Repo, __p.Ver).Scan(&_pid)
	var _pe *pgconn.PgError
	if errors.As(_err, &_pe) && _pe.Code == "23505" {
		return 0, 0, ErrSlug
	}
	if _err != nil {
		return 0, 0, _err
	}
	__q.Plug, __q.Kind = _pid, "new"
	if _qid, _err = add(__ctx, _tx, __q, __pic); _err != nil {
		return 0, 0, _err
	}
	return _pid, _qid, _tx.Commit(__ctx)
}

type _Q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func add(__ctx context.Context, __db _Q, __q Req, __pic []byte) (int64, error) {
	var _id int64
	return _id, __db.QueryRow(__ctx, `
		insert into reqs (user_id, plug_id, kind, name, about, ver, note, repo, pic) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`,
		__q.User, __q.Plug, __q.Kind, __q.Name, __q.About, __q.Ver, __q.Note, __q.Repo, __pic).Scan(&_id)
}

func (r *Repo) Add(__ctx context.Context, __q Req, __pic []byte) (int64, error) {
	return add(__ctx, r.__db, __q, __pic)
}

func (r *Repo) Busy(__ctx context.Context, __plug int64) (bool, error) {
	var _n int
	return _n > 0, r.__db.QueryRow(__ctx, `select count(*) from reqs where plug_id = $1 and status = 'wait'`, __plug).Scan(&_n)
}

func (r *Repo) History(__ctx context.Context, __plug int64) ([]Req, error) {
	return r.reqs(__ctx, `select `+__cols__+` from reqs where plug_id = $1 order by created desc limit 100`, __plug)
}

func (r *Repo) Apply(__ctx context.Context, __q *Req) error {
	_tx, _err := r.__db.Begin(__ctx)
	if _err != nil {
		return _err
	}
	defer _tx.Rollback(__ctx)
	switch __q.Kind {
	case "edit":
		_, _err = _tx.Exec(__ctx, `update plugs set name = $2, about = $3, updated = now() where id = $1`, __q.Plug, __q.Name, __q.About)
	default:
		_, _err = _tx.Exec(__ctx, `update plugs set name = $2, about = $3, ver = $4, repo = $5, live = true, updated = now() where id = $1`,
			__q.Plug, __q.Name, __q.About, __q.Ver, __q.Repo)
	}
	if _err != nil {
		return _err
	}
	if __q.Pic {
		if _, _err = _tx.Exec(__ctx, `update plugs set pic = (select pic from reqs where id = $2), pv = pv + 1 where id = $1`, __q.Plug, __q.ID); _err != nil {
			return _err
		}
	}
	return _tx.Commit(__ctx)
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `
		with p as (delete from plugs where user_id = $1 returning id),
		q as (delete from reqs where user_id = $1 returning id)
		delete from devs where user_id = $1`, __user)
	return _err
}
