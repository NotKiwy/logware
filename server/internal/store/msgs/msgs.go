package msgs

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Row struct {
	ID      int64
	User    int64
	Conn    string
	Chat    int64
	Msg     int64
	Ver     int
	Out     bool
	Kind    string
	Body    []byte
	Has     bool
	Created time.Time
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

const __cols__ = `id, user_id, conn_id, chat_id, msg_id, ver, out, kind, body, media is not null, created`

func (r *Repo) Add(__ctx context.Context, __r *Row) (int64, error) {
	var _id int64
	_err := r.__db.QueryRow(__ctx, `
		insert into msgs (user_id, conn_id, chat_id, msg_id, ver, out, kind, body)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (conn_id, chat_id, msg_id, ver) do nothing
		returning id`,
		__r.User, __r.Conn, __r.Chat, __r.Msg, __r.Ver, __r.Out, __r.Kind, __r.Body).Scan(&_id)
	if errors.Is(_err, pgx.ErrNoRows) {
		return 0, nil
	}
	return _id, _err
}

func (r *Repo) SetMedia(__ctx context.Context, __id int64, __media []byte) error {
	_, _err := r.__db.Exec(__ctx, `update msgs set media = $2 where id = $1`, __id, __media)
	return _err
}

func (r *Repo) Last(__ctx context.Context, __conn string, __chat, __msg int64) (*Row, error) {
	_rows, _err := r.many(__ctx, `select `+__cols__+` from msgs where conn_id = $1 and chat_id = $2 and msg_id = $3 order by ver desc limit 1`, __conn, __chat, __msg)
	if _err != nil || len(_rows) == 0 {
		return nil, _err
	}
	return &_rows[0], nil
}

func (r *Repo) Find(__ctx context.Context, __conn string, __chat int64, __ids []int64) ([]Row, error) {
	return r.many(__ctx, `select `+__cols__+` from msgs where conn_id = $1 and chat_id = $2 and msg_id = any($3) order by msg_id, ver`, __conn, __chat, __ids)
}

type Chat struct {
	ID   int64     `json:"id"`
	Body []byte    `json:"body"`
	N    int       `json:"n"`
	Last time.Time `json:"-"`
}

func (r *Repo) Chats(__ctx context.Context, __user int64) ([]Chat, error) {
	_rs, _err := r.__db.Query(__ctx, `
		select c.chat_id, c.n, c.last, m.body from (
			select chat_id, count(*) n, max(id) top, max(created) last from msgs where user_id = $1 group by chat_id
		) c join msgs m on m.id = c.top order by c.last desc limit 300`, __user)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	_out := []Chat{}
	for _rs.Next() {
		var _c Chat
		if _err = _rs.Scan(&_c.ID, &_c.N, &_c.Last, &_c.Body); _err != nil {
			return nil, _err
		}
		_out = append(_out, _c)
	}
	return _out, _rs.Err()
}

func (r *Repo) Page(__ctx context.Context, __user int64, __ids []int64, __off, __lim int) ([]Row, error) {
	return r.many(__ctx, `select `+__cols__+` from msgs where user_id = $1 and id = any($2) order by msg_id, ver offset $3 limit $4`, __user, __ids, __off, __lim)
}

func (r *Repo) Media(__ctx context.Context, __user, __id int64) ([]byte, error) {
	var _b []byte
	_err := r.__db.QueryRow(__ctx, `select media from msgs where user_id = $1 and id = $2`, __user, __id).Scan(&_b)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	return _b, _err
}

func (r *Repo) DropOne(__ctx context.Context, __conn string, __chat, __msg int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from msgs where conn_id = $1 and chat_id = $2 and msg_id = $3`, __conn, __chat, __msg)
	return _err
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from msgs where user_id = $1`, __user)
	return _err
}

func (r *Repo) Clean(__ctx context.Context) (int64, error) {
	_tag, _err := r.__db.Exec(__ctx, `
		delete from msgs m using users u
		where m.user_id = u.id
			and coalesce((u.prefs->>'keep')::int, 30) > 0
			and m.created < now() - make_interval(days => coalesce((u.prefs->>'keep')::int, 30))`)
	return _tag.RowsAffected(), _err
}

func (r *Repo) many(__ctx context.Context, __q string, __args ...any) ([]Row, error) {
	_rs, _err := r.__db.Query(__ctx, __q, __args...)
	if _err != nil {
		return nil, _err
	}
	defer _rs.Close()
	var _out []Row
	for _rs.Next() {
		var _r Row
		if _err = _rs.Scan(&_r.ID, &_r.User, &_r.Conn, &_r.Chat, &_r.Msg, &_r.Ver, &_r.Out, &_r.Kind, &_r.Body, &_r.Has, &_r.Created); _err != nil {
			return nil, _err
		}
		_out = append(_out, _r)
	}
	return _out, _rs.Err()
}

func (r *Repo) Count(__ctx context.Context, __user int64) (int, error) {
	var _n int
	_err := r.__db.QueryRow(__ctx, `select count(distinct (conn_id, chat_id, msg_id)) from msgs where user_id = $1`, __user).Scan(&_n)
	return _n, _err
}

func (r *Repo) HasChat(__ctx context.Context, __user, __chat int64) (bool, error) {
	var _ok bool
	_err := r.__db.QueryRow(__ctx, `select exists (select 1 from msgs where user_id = $1 and chat_id = $2)`, __user, __chat).Scan(&_ok)
	return _ok, _err
}
