package once

import (
	"context"

	"logware/internal/bot/keep"
	"logware/internal/bot/notify"
	"logware/internal/bot/own"
	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/types"
)

func Cmd(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message, _ []string) error {
	if __m.Reply == nil {
		return nil
	}
	if _f, _ := __m.Reply.Media(); _f == nil {
		return nil
	}
	if _err := grab(__ctx, __d, __o, __m.Reply); _err != nil {
		return _err
	}
	if !__o.User.Prefs.OnceHide || !__o.Conn.Rights.DelAll {
		return nil
	}
	if _err := __d.Msgs.DropOne(__ctx, __o.Conn.ID, __m.Chat.ID, __m.ID); _err != nil {
		return _err
	}
	return __d.Tg.DelBiz(__ctx, __o.Conn.ID, []int64{__m.ID})
}

func Auto(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message) error {
	_r := __m.Reply
	if _f, _ := _r.Media(); _f == nil {
		return nil
	}
	_prev, _err := __d.Msgs.Last(__ctx, __o.Conn.ID, _r.Chat.ID, _r.ID)
	if _err != nil || _prev == nil || (_prev.Kind != "empty" && _prev.Kind != "other") {
		return _err
	}
	return grab(__ctx, __d, __o, _r)
}

func grab(__ctx context.Context, __d *deps.D, __o *own.Own, __r *types.Message) error {
	_, _ver, _err := keep.Next(__ctx, __d, __o.Conn.ID, __r.Chat.ID, __r.ID)
	if _err != nil {
		return _err
	}
	_id, _err := keep.Msg(__ctx, __d, __o, __r, _ver, true)
	if _err != nil || _id == 0 {
		return _err
	}
	return notify.Send(__ctx, __d, __o, "once", []int64{_id}, text.Saved(text.Name(__r.Chat)))
}
