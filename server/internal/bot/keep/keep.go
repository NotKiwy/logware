package keep

import (
	"context"
	"log/slog"

	"logware/internal/bot/grab"
	"logware/internal/bot/own"
	"logware/internal/crypto/seal"
	"logware/internal/deps"
	"logware/internal/store/msgs"
	"logware/internal/tg/types"
)

func Msg(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message, __ver int, __force bool) (int64, error) {
	_body, _err := seal.Seal(__o.User.Pub, __m.Raw)
	if _err != nil {
		return 0, _err
	}
	_id, _err := __d.Msgs.Add(__ctx, &msgs.Row{
		User: __o.User.ID,
		Conn: __o.Conn.ID,
		Chat: __m.Chat.ID,
		Msg:  __m.ID,
		Ver:  __ver,
		Out:  __m.Out(__o.Conn.User),
		Kind: __m.Kind(),
		Body: _body,
	})
	if _err != nil || _id == 0 {
		return _id, _err
	}
	_f, _k := __m.Media()
	if _f == nil || !(__force || __o.User.Prefs.Wants(_k)) {
		return _id, nil
	}
	_b, _err := grab.Media(__ctx, __d, __o.User.Pub, _f)
	if _err != nil {
		slog.Warn("media skipped", "row", _id, "err", _err)
		return _id, nil
	}
	return _id, __d.Msgs.SetMedia(__ctx, _id, _b)
}

func Next(__ctx context.Context, __d *deps.D, __conn string, __chat, __msg int64) (*msgs.Row, int, error) {
	_prev, _err := __d.Msgs.Last(__ctx, __conn, __chat, __msg)
	if _err != nil || _prev == nil {
		return nil, 0, _err
	}
	return _prev, _prev.Ver + 1, nil
}
