package boot

import (
	"context"

	"logware/internal/cfg"
	"logware/internal/deps"
	"logware/internal/store/conns"
	"logware/internal/store/db"
	"logware/internal/store/devs"
	"logware/internal/store/events"
	"logware/internal/store/ign"
	"logware/internal/store/kws"
	"logware/internal/store/mods"
	"logware/internal/store/msgs"
	"logware/internal/store/promos"
	"logware/internal/store/subs"
	"logware/internal/store/users"
	"logware/internal/tg/client"
)

func Open(__ctx context.Context, __c *cfg.Cfg) (*deps.D, func(), error) {
	_pool, _err := db.Open(__ctx, __c.DB)
	if _err != nil {
		return nil, nil, _err
	}
	_d := &deps.D{
		DB:     _pool,
		Cfg:    __c,
		Tg:     client.New(__c.Tok),
		Users:  users.New(_pool),
		Conns:  conns.New(_pool),
		Msgs:   msgs.New(_pool),
		Evs:    events.New(_pool),
		Mods:   mods.New(_pool),
		Ign:    ign.New(_pool),
		Kws:    kws.New(_pool),
		Subs:   subs.New(_pool),
		Devs:   devs.New(_pool),
		Promos: promos.New(_pool),
	}
	return _d, _pool.Close, nil
}
