package start

import (
	"context"
	"strings"

	"logware/internal/bot/help"
	"logware/internal/bot/review"
	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

func Run(__ctx context.Context, __d *deps.D, __m *types.Message) error {
	if __m.Chat.Type != "private" || __m.From == nil || __m.From.Bot {
		return nil
	}
	if _done, _err := review.Why(__ctx, __d, __m); _done || _err != nil {
		return _err
	}
	if _done, _err := help.Run(__ctx, __d, __m); _done || _err != nil {
		return _err
	}
	_, _ = __d.Users.Touch(__ctx, __m.From.ID, strings.TrimSpace(__m.From.First+" "+__m.From.Last), __m.From.Nick, __m.From.Prem)
	_u, _err := __d.Users.Get(__ctx, __m.From.ID)
	if _err != nil {
		return _err
	}
	if _u != nil && _u.Banned {
		return __d.Tg.Send(__ctx, __m.Chat.ID, text.Banned(), nil)
	}
	if !_u.Ready() {
		return __d.Tg.Send(__ctx, __m.Chat.ID, text.Hello(), kb.One(kb.Web(text.RegBtn, __d.Url("/reg")).With(text.Logo)))
	}
	_c, _err := __d.Conns.ByUser(__ctx, __m.From.ID)
	if _err != nil {
		return _err
	}
	_on := _c != nil && _c.Active
	_kb := kb.One(kb.Web(text.SetBtn, __d.Url("/")).With(text.Gear))
	if !_on {
		_kb = kb.Rows([]kb.Btn{kb.Url(text.ProfBtn, text.ProfURL).With(text.Gear)}, []kb.Btn{kb.Web(text.SetBtn, __d.Url("/")).With(text.Logo)})
	}
	return __d.Tg.Send(__ctx, __m.Chat.ID, text.Ready(_on, __d.Nick), _kb)
}
