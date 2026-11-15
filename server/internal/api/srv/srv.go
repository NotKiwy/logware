package srv

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"logware/internal/api/auth"
	"logware/internal/api/ava"
	"logware/internal/api/chats"
	"logware/internal/api/dev"
	"logware/internal/api/ev"
	"logware/internal/api/hist"
	"logware/internal/api/hook"
	"logware/internal/api/kwapi"
	"logware/internal/api/limit"
	"logware/internal/api/me"
	"logware/internal/api/mods"
	"logware/internal/api/pay"
	"logware/internal/api/prefs"
	"logware/internal/api/promo"
	"logware/internal/api/reg"
	"logware/internal/api/send"
	"logware/internal/api/web"
	"logware/internal/api/wipe"
	"logware/internal/deps"
)

func Run(__ctx context.Context, __d *deps.D) error {
	_a := auth.Mw(__d.Cfg.Tok, func(__c context.Context, __w auth.Who) bool {
		_b, _err := __d.Users.Touch(__c, __w.ID, __w.Full(), __w.Nick, __w.Prem)
		if _err != nil {
			slog.Warn("touch", "err", _err)
		}
		return !_b
	})
	_uid := func(__p string) func(*http.Request) string {
		return func(r *http.Request) string { return __p + strconv.FormatInt(auth.ID(r), 10) }
	}
	_lim := func(__rate, __burst float64, __p string) func(http.Handler) http.Handler {
		_l := limit.New(__rate, __burst)
		return func(__h http.Handler) http.Handler { return _a(limit.Mw(_l, _uid(__p))(__h)) }
	}
	_u := _lim(10, 60, "u")
	_med := _lim(25, 150, "m")
	_slow := _lim(0.5, 5, "s")
	_rare := _lim(0.1, 3, "r")
	_ip := limit.Mw(limit.New(20, 80), limit.IP)
	_hk := limit.Mw(limit.New(10, 30), limit.IP)
	_m := http.NewServeMux()
	_m.HandleFunc("GET /api/bot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		_, _ = w.Write([]byte(`{"bot":"` + __d.Nick + `"}`))
	})
	_m.Handle("GET /api/me", _u(me.Get(__d)))
	_m.Handle("POST /api/reg", _rare(reg.Post(__d)))
	_m.Handle("POST /api/prefs", _u(prefs.Post(__d)))
	_m.Handle("POST /api/wipe", _rare(wipe.Post(__d)))
	_m.Handle("GET /api/ev/{id}", _u(ev.Get(__d)))
	_m.Handle("GET /api/media/{id}", _med(ev.Media(__d)))
	_m.Handle("POST /api/send", _slow(send.Post(__d)))
	_m.Handle("GET /api/ava/{id}", _med(ava.Get(__d)))
	_m.Handle("GET /api/hist", _u(hist.Get(__d)))
	_m.Handle("POST /api/clear", _u(hist.Clear(__d)))
	_m.Handle("GET /api/chats", _u(chats.Get(__d)))
	_m.Handle("POST /api/chats", _u(chats.Post(__d)))
	_m.Handle("GET /api/kws", _u(kwapi.Get(__d)))
	_m.Handle("POST /api/kws", _u(kwapi.Post(__d)))
	_m.Handle("POST /api/kws/del", _u(kwapi.Del(__d)))
	_m.Handle("GET /api/mods", _u(mods.Get(__d)))
	_m.Handle("POST /api/mods/{id}", _u(mods.Post(__d)))
	_m.Handle("POST /api/pay", _slow(pay.Post(__d)))
	_m.Handle("POST /api/promo", _lim(0.05, 5, "p")(promo.Post(__d)))
	_m.Handle("GET /api/dev", _u(dev.Get(__d)))
	_m.Handle("POST /api/dev", _u(dev.Post(__d)))
	_m.Handle("POST /api/dev/plug", _slow(dev.New(__d)))
	_m.Handle("GET /api/dev/plug/{id}", _u(dev.Show(__d)))
	_m.Handle("POST /api/dev/plug/{id}", _slow(dev.Change(__d)))
	_m.Handle("GET /api/dev/req/{id}/pic", _med(dev.ReqPic(__d)))
	_m.Handle("GET /pic/{id}", mods.Pic(__d))
	_m.Handle("POST /hook", _hk(hook.Any(__d)))
	_m.Handle("POST /hook/cb", _hk(hook.Cb(__d)))
	_m.Handle("POST /hook/xr", _hk(hook.Xr(__d)))
	_m.Handle("/api/", http.NotFoundHandler())
	_m.Handle("/", web.Static(__d.Cfg.Web))
	_s := &http.Server{Addr: __d.Cfg.Addr, Handler: limit.Hide(_ip(_m)), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-__ctx.Done()
		_c, _done := context.WithTimeout(context.Background(), 5*time.Second)
		defer _done()
		_ = _s.Shutdown(_c)
	}()
	slog.Info("http", "addr", __d.Cfg.Addr)
	if _err := _s.ListenAndServe(); !errors.Is(_err, http.ErrServerClosed) {
		return _err
	}
	return nil
}
