package conns

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"logware/internal/tg/types"
)

type Conn struct {
	ID     string
	User   int64
	Chat   int64
	Rights types.Rights
	Active bool
}

type Repo struct {
	__db *pgxpool.Pool
}

func New(__db *pgxpool.Pool) *Repo {
	return &Repo{__db: __db}
}

func (r *Repo) Put(__ctx context.Context, __c *types.Conn) error {
	_rights, _err := json.Marshal(__c.Rights)
	if _err != nil {
		return _err
	}
	_, _err = r.__db.Exec(__ctx, `
		insert into conns (id, user_id, chat_id, rights, active, name, nick) values ($1, $2, $3, $4, $5, $6, $7)
		on conflict (id) do update set user_id = excluded.user_id, chat_id = excluded.chat_id,
			rights = excluded.rights, active = excluded.active, updated = now(),
			name = coalesce(nullif(excluded.name, ''), conns.name), nick = coalesce(nullif(excluded.nick, ''), conns.nick)`,
		__c.ID, __c.User.ID, __c.ChatID, _rights, __c.On, strings.TrimSpace(__c.User.First+" "+__c.User.Last), __c.User.Nick)
	return _err
}

func (r *Repo) Get(__ctx context.Context, __id string) (*Conn, error) {
	return r.one(__ctx, `select id, user_id, chat_id, rights, active from conns where id = $1`, __id)
}

func (r *Repo) ByUser(__ctx context.Context, __user int64) (*Conn, error) {
	return r.one(__ctx, `select id, user_id, chat_id, rights, active from conns where user_id = $1 order by active desc, updated desc limit 1`, __user)
}

func (r *Repo) one(__ctx context.Context, __q string, __arg any) (*Conn, error) {
	_c := &Conn{}
	var _raw []byte
	_err := r.__db.QueryRow(__ctx, __q, __arg).Scan(&_c.ID, &_c.User, &_c.Chat, &_raw, &_c.Active)
	if errors.Is(_err, pgx.ErrNoRows) {
		return nil, nil
	}
	if _err != nil {
		return nil, _err
	}
	_ = json.Unmarshal(_raw, &_c.Rights)
	return _c, nil
}

func (r *Repo) Drop(__ctx context.Context, __user int64) error {
	_, _err := r.__db.Exec(__ctx, `delete from conns where user_id = $1`, __user)
	return _err
}
