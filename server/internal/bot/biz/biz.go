package biz

import (
	"context"

	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

func Run(__ctx context.Context, __d *deps.D, __c *types.Conn) error {
	if _err := __d.Conns.Put(__ctx, __c); _err != nil {
		return _err
	}
	if !__c.On {
		return __d.Tg.Send(__ctx, __c.ChatID, text.Unlinked(), nil)
	}
	_u, _err := __d.Users.Get(__ctx, __c.User.ID)
	if _err != nil {
		return _err
	}
	var _kb *kb.Kb
	if !_u.Ready() {
		_kb = kb.One(kb.Web(text.RegBtn, __d.Url("/reg")).With(text.Logo))
	} else if !__c.Rights.Reply || !__c.Rights.DelAll {
		_kb = kb.One(kb.Url(text.ProfBtn, text.ProfURL).With(text.Gear))
	}
	return __d.Tg.Send(__ctx, __c.ChatID, text.Linked(__c.Rights, _u.Ready(), __d.Nick), _kb)
}
