package own

import (
	"context"

	"logware/internal/deps"
	"logware/internal/store/conns"
	"logware/internal/store/users"
)

type Own struct {
	Conn *conns.Conn
	User *users.User
}

func (o *Own) Ready() bool {
	return o != nil && o.Conn != nil && o.User.Live()
}

func Get(__ctx context.Context, __d *deps.D, __id string) (*Own, error) {
	_c, _err := __d.Conns.Get(__ctx, __id)
	if _err != nil {
		return nil, _err
	}
	if _c == nil {
		_tc, _err := __d.Tg.Conn(__ctx, __id)
		if _err != nil {
			return nil, _err
		}
		if _err = __d.Conns.Put(__ctx, _tc); _err != nil {
			return nil, _err
		}
		if _c, _err = __d.Conns.Get(__ctx, __id); _err != nil {
			return nil, _err
		}
	}
	_u, _err := __d.Users.Get(__ctx, _c.User)
	if _err != nil {
		return nil, _err
	}
	return &Own{Conn: _c, User: _u}, nil
}
