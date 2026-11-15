package app

import (
	"context"
	"log/slog"

	"logware/internal/api/srv"
	"logware/internal/boot"
	"logware/internal/bot/pay"
	"logware/internal/bot/poll"
	"logware/internal/bot/text"
	"logware/internal/cfg"
	"logware/internal/jobs/clean"
	"logware/internal/mods/clock"
	"logware/internal/tg/client"
	"logware/internal/tunnel"
)

func Run(__ctx context.Context) error {
	_c, _err := cfg.Load()
	if _err != nil {
		return _err
	}
	if _c.Pub == "" {
		if _c.Pub, _err = tunnel.Find(__ctx); _err != nil {
			return _err
		}
	}
	_d, _done, _err := boot.Open(__ctx, _c)
	if _err != nil {
		return _err
	}
	defer _done()
	_me, _err := _d.Tg.Me(__ctx)
	if _err != nil {
		return _err
	}
	_d.Nick = _me.Nick
	if _c.SupT != "" {
		if _s, _err := client.New(_c.SupT).Me(__ctx); _err == nil {
			_d.Sup = _s.Nick
		}
	}
	if !_me.Biz {
		slog.Warn("business mode is off, enable Secretary Mode for the bot in @BotFather")
	}
	if _err = _d.Tg.NoHook(__ctx); _err != nil {
		return _err
	}
	if _err = _d.Tg.Cmds(__ctx, [][2]string{{"start", "Открыть Logware"}, {"paysupport", "Поддержка по оплате"}, {"terms", "Пользовательское соглашение"}, {"privacy", "Политика конфиденциальности"}}); _err != nil {
		slog.Warn("commands", "err", _err)
	}
	if _err = _d.Tg.Menu(__ctx, text.MenuBtn, _d.Url("/")); _err != nil {
		slog.Warn("menu button", "err", _err)
	}
	slog.Info("started", "bot", "@"+_me.Nick, "url", _c.Pub)
	_errs := make(chan error, 1)
	go func() { _errs <- srv.Run(__ctx, _d) }()
	go clean.Run(__ctx, _d)
	go clock.Run(__ctx, _d)
	go pay.Remind(__ctx, _d)
	go pay.Sweep(__ctx, _d)
	go func() { _errs <- poll.Run(__ctx, _d) }()
	return <-_errs
}
