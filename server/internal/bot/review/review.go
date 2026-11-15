package review

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/store/devs"
	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

func Ask(__ctx context.Context, __d *deps.D, __dev *devs.Dev, __q devs.Req, __slug string) {
	if __d.Cfg.Own == 0 {
		return
	}
	_id := strconv.FormatInt(__q.ID, 10)
	_kb := kb.Rows(
		[]kb.Btn{kb.Url("Открыть код", __q.Repo).With(text.Eye)},
		[]kb.Btn{kb.Data("Одобрить", "rq:ok:"+_id).With(text.Check), kb.Data("Отклонить", "rq:no:"+_id).With(text.Alert)},
	)
	_msg, _err := __d.Tg.Post(__ctx, __d.Cfg.Own, text.Review(__dev.Nick, __q.User, __q.Kind, __slug, __q.Name, __q.Ver, __q.About, __q.Note, __q.Repo, __q.Pic), _kb)
	if _err != nil {
		slog.Error("review: send", "err", _err)
		return
	}
	_ = __d.Devs.Msg(__ctx, __q.ID, _msg)
}

func Decide(__ctx context.Context, __d *deps.D, __q *devs.Req, __ok bool, __why string) error {
	if __q.Status != "wait" {
		return nil
	}
	_st := "no"
	if __ok {
		_st, __why = "ok", ""
	}
	if __ok && __q.Plug != 0 {
		if _err := __d.Devs.Apply(__ctx, __q); _err != nil {
			return _err
		}
	}
	if _err := __d.Devs.Mark(__ctx, __q.ID, _st, __why); _err != nil {
		return _err
	}
	if __q.Msg != 0 && __d.Cfg.Own != 0 {
		_ = __d.Tg.Retext(__ctx, __d.Cfg.Own, __q.Msg, text.Reviewed(__q.Kind, __q.Name, __q.Ver, __q.Repo, __ok, __why))
	}
	return __d.Tg.Send(__ctx, __q.User, text.Verdict(__q.Kind, __q.Name, __q.Ver, __ok, __why), kb.One(kb.Web(text.DevBtn, __d.Url("/dev")).With(text.Gear)))
}

func Cb(__ctx context.Context, __d *deps.D, __c *types.Callback) error {
	_p := strings.Split(__c.Data, ":")
	if len(_p) != 3 || _p[0] != "rq" {
		return nil
	}
	if __c.From.ID != __d.Cfg.Own || __d.Cfg.Own == 0 {
		return __d.Tg.Answer(__ctx, __c.ID, "")
	}
	_id, _ := strconv.ParseInt(_p[2], 10, 64)
	_q, _err := __d.Devs.Req(__ctx, _id)
	if _err != nil || _q == nil {
		return __d.Tg.Answer(__ctx, __c.ID, "Заявка не найдена")
	}
	if _p[1] != "ok" && _p[1] != "no" {
		return __d.Tg.Answer(__ctx, __c.ID, "")
	}
	_ = __d.Tg.Answer(__ctx, __c.ID, "Готово")
	return Decide(__ctx, __d, _q, _p[1] == "ok", "")
}

func Why(__ctx context.Context, __d *deps.D, __m *types.Message) (bool, error) {
	if __d.Cfg.Own == 0 || __m.From == nil || __m.From.ID != __d.Cfg.Own || __m.Reply == nil || __m.Text == "" {
		return false, nil
	}
	_q, _err := __d.Devs.ByMsg(__ctx, __m.Reply.ID)
	if _err != nil || _q == nil {
		return false, _err
	}
	if _err = Decide(__ctx, __d, _q, false, __m.Text); _err != nil {
		return true, _err
	}
	return true, __d.Tg.Send(__ctx, __m.Chat.ID, text.WhySent(), nil)
}
