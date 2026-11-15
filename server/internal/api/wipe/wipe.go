package wipe

import (
	"context"
	"log/slog"
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/mods/reg"
	"logware/internal/pay"
)

func off(__ctx context.Context, __d *deps.D, __id int64) {
	_st, _err := __d.Mods.List(__ctx, __id)
	if _err != nil {
		return
	}
	for _, _m := range reg.All() {
		if _s, _ok := _st[_m.ID]; _ok && _s.On && _m.Set != nil {
			if _, _err = _m.Set(__ctx, __d, __id, false, _m.Fill(_s.Cfg)); _err != nil {
				slog.Warn("wipe: mod off", "mod", _m.ID, "err", _err)
			}
		}
	}
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id, _ctx := auth.ID(r), r.Context()
		if _s, _err := __d.Subs.Get(_ctx, _id); _err == nil && _s.Live() && pay.Star(_s.Charge) {
			if _err = __d.Tg.Cancel(_ctx, _id, _s.Charge); _err != nil {
				slog.Error("wipe: cancel sub", "user", _id, "err", _err)
			}
		}
		off(_ctx, __d, _id)
		for _, _fn := range []func() error{
			func() error { return __d.Msgs.Drop(_ctx, _id) },
			func() error { return __d.Evs.Drop(_ctx, _id) },
			func() error { return __d.Mods.Drop(_ctx, _id) },
			func() error { return __d.Ign.Drop(_ctx, _id) },
			func() error { return __d.Kws.Drop(_ctx, _id) },
			func() error { return __d.Devs.Drop(_ctx, _id) },
			func() error { return __d.Conns.Drop(_ctx, _id) },
			func() error { _, _err := __d.Subs.Take(_ctx, _id); return _err },
			func() error { return __d.Users.Drop(_ctx, _id) },
		} {
			if _err := _fn(); _err != nil {
				out.Oops(w, _err)
				return
			}
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
