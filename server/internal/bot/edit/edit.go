package edit

import (
	"context"

	"logware/internal/bot/keep"
	"logware/internal/bot/notify"
	"logware/internal/bot/own"
	"logware/internal/bot/text"
	"logware/internal/deps"
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
	if _skip, _err := __d.Ign.Has(__ctx, _o.User.ID, __m.Chat.ID); _err != nil || _skip {
		return _err
	}
	_prev, _ver, _err := keep.Next(__ctx, __d, _o.Conn.ID, __m.Chat.ID, __m.ID)
	if _err != nil {
		return _err
	}
	_id, _err := keep.Msg(__ctx, __d, _o, __m, _ver, false)
	if _err != nil || _id == 0 || _prev == nil {
		return _err
	}
	_p := _o.User.Prefs
	if !_p.Edit || (_prev.Out && !_p.EditMe) {
		return nil
	}
	return notify.Send(__ctx, __d, _o, "edit", []int64{_prev.ID, _id}, text.Edit(text.Name(__m.Chat)))
}
