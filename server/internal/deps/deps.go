package deps

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"logware/internal/cfg"
	"logware/internal/store/conns"
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

type D struct {
	DB     *pgxpool.Pool
	Cfg    *cfg.Cfg
	Nick   string
	Sup    string
	Tg     *client.Client
	Users  *users.Repo
	Conns  *conns.Repo
	Msgs   *msgs.Repo
	Evs    *events.Repo
	Mods   *mods.Repo
	Ign    *ign.Repo
	Kws    *kws.Repo
	Subs   *subs.Repo
	Devs   *devs.Repo
	Promos *promos.Repo
}

func (d *D) Url(__path string) string {
	return d.Cfg.Pub + __path
}

func (d *D) Pro(__ctx context.Context, __user int64) bool {
	if __user != 0 && __user == d.Cfg.Own {
		return true
	}
	_s, _err := d.Subs.Get(__ctx, __user)
	return _err == nil && _s.Live() && (_s.Plan == "pro" || _s.Plan == "dev")
}
