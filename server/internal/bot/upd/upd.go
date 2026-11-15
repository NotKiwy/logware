package upd

import (
	"context"

	"logware/internal/bot/biz"
	"logware/internal/bot/del"
	"logware/internal/bot/edit"
	"logware/internal/bot/pay"
	"logware/internal/bot/review"
	"logware/internal/bot/save"
	"logware/internal/bot/start"
	"logware/internal/deps"
	"logware/internal/tg/types"
)

func Run(__ctx context.Context, __d *deps.D, __u *types.Update) error {
	switch {
	case __u.Cb != nil:
		return review.Cb(__ctx, __d, __u.Cb)
	case __u.Pre != nil:
		return pay.Pre(__ctx, __d, __u.Pre)
	case __u.Msg != nil && __u.Msg.Paid != nil:
		return pay.Done(__ctx, __d, __u.Msg)
	case __u.Msg != nil:
		return start.Run(__ctx, __d, __u.Msg)
	case __u.Conn != nil:
		return biz.Run(__ctx, __d, __u.Conn)
	case __u.BizMsg != nil:
		return save.Run(__ctx, __d, __u.BizMsg)
	case __u.BizEdit != nil:
		return edit.Run(__ctx, __d, __u.BizEdit)
	case __u.BizDel != nil:
		return del.Run(__ctx, __d, __u.BizDel)
	}
	return nil
}
