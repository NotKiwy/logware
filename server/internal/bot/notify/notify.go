package notify

import (
	"context"

	"logware/internal/bot/own"
	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/kb"
)

func Send(__ctx context.Context, __d *deps.D, __o *own.Own, __kind string, __refs []int64, __msg string) error {
	return push(__ctx, __d, __o, __kind, __refs, __msg, __o.User.Prefs.Mute(__kind))
}

func Loud(__ctx context.Context, __d *deps.D, __o *own.Own, __kind string, __refs []int64, __msg string) error {
	return push(__ctx, __d, __o, __kind, __refs, __msg, false)
}

func push(__ctx context.Context, __d *deps.D, __o *own.Own, __kind string, __refs []int64, __msg string, __mute bool) error {
	_id, _err := __d.Evs.Add(__ctx, __o.User.ID, __kind, __refs)
	if _err != nil {
		return _err
	}
	return __d.Tg.Note(__ctx, __o.Conn.Chat, __msg, kb.One(kb.Web(text.ShowBtn, __d.Url("/e/"+_id)).With(text.Eye)), __mute)
}
