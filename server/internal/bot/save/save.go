package save

import (
	"context"
	"log/slog"

	"logware/internal/bot/keep"
	"logware/internal/bot/kw"
	"logware/internal/bot/once"
	"logware/internal/bot/own"
	"logware/internal/bot/route"
	"logware/internal/deps"
	"logware/internal/mods/reg"
	"logware/internal/tg/types"
)

func Run(__ctx context.Context, __d *deps.D, __m *types.Message) error {
	if __m.ByBot != nil {
		return nil
	}
	_o, _err := own.Get(__ctx, __d, __m.Conn)
	if _err != nil || !_o.Ready() {
		return _err
	}
	_out := __m.Out(_o.Conn.User)
	_skip, _err := __d.Ign.Has(__ctx, _o.User.ID, __m.Chat.ID)
	if _err != nil {
		return _err
	}
	if !_skip {
		_row, _err := keep.Msg(__ctx, __d, _o, __m, 0, false)
		if _err != nil {
			return _err
		}
		if !_out {
			if _err = kw.Check(__ctx, __d, _o, __m, _row); _err != nil {
				slog.Warn("kw", "err", _err)
			}
		}
		reg.Msg(__ctx, __d, _o, __m)
	}
	if !_out {
		return nil
	}
	if route.IsCmd(__m.Text) {
		return route.Run(__ctx, __d, _o, __m)
	}
	if __m.Reply != nil && _o.User.Prefs.Once && _o.User.Prefs.OnceAuto {
		return once.Auto(__ctx, __d, _o, __m)
	}
	return nil
}
