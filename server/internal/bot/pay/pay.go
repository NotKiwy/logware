package pay

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"logware/internal/bot/text"
	"logware/internal/coin/cb"
	"logware/internal/coin/xr"
	"logware/internal/deps"
	"logware/internal/pay"
	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

func Pre(__ctx context.Context, __d *deps.D, __q *types.Pre) error {
	_p, _u, _ok := pay.Parse(__q.Payload)
	if !_ok || _u != __q.From.ID || __q.Currency != "XTR" || __q.Total != _p.Price {
		return __d.Tg.PreOk(__ctx, __q.ID, "Счёт устарел, открой оплату в приложении заново.")
	}
	return __d.Tg.PreOk(__ctx, __q.ID, "")
}

func Done(__ctx context.Context, __d *deps.D, __m *types.Message) error {
	_x := __m.Paid
	_p, _u, _ok := pay.Parse(_x.Payload)
	if !_ok || __m.From == nil || _u != __m.From.ID {
		slog.Error("pay: bad payload", "payload", _x.Payload, "charge", _x.Charge)
		return nil
	}
	_until := time.Unix(_x.Expires, 0)
	if _x.Expires == 0 {
		_until = time.Now().Add(pay.Period * time.Second)
	}
	_old, _ := __d.Subs.Get(__ctx, _u)
	_plan := _p.ID
	if _old.Live() {
		_plan, _until = pay.Merge(_old.Plan, _old.Until, _p.ID, pay.Period*time.Second, time.Now())
	}
	_new, _err := __d.Subs.Pay(__ctx, _u, _plan, _x.Total, "XTR", _until, _x.Charge)
	if _err != nil || !_new {
		return _err
	}
	_rec := _x.Recurring
	if _plan != _p.ID {
		_rec = false
		slog.Warn("pay: lower plan renewal", "user", _u, "plan", _p.ID, "live", _plan, "charge", _x.Charge)
		if _x.Recurring {
			_ = __d.Tg.Cancel(__ctx, _u, _x.Charge)
		}
		if _q, _ok := pay.Get(_plan); _ok {
			_p = _q
		}
	} else if _old.Live() && _old.Plan != _p.ID && pay.Star(_old.Charge) {
		if _err = __d.Tg.Cancel(__ctx, _u, _old.Charge); _err != nil {
			slog.Error("pay: cancel old sub", "user", _u, "err", _err)
		}
	}
	slog.Info("pay", "user", _u, "plan", _p.ID, "stars", _x.Total, "until", _until)
	return __d.Tg.Send(__ctx, __m.Chat.ID, text.Paid(_p.Name, _until, _rec), kb.One(kb.Web(text.SetBtn, __d.Url("/")).With(text.Gear)))
}

func Coin(__ctx context.Context, __d *deps.D, __payload, __charge string) error {
	if _, _t := pay.Tag(__payload); _t != __d.Cfg.Tag && (_t != "" || __d.Cfg.Tag == "app") {
		return nil
	}
	_p, _u, _ok := pay.Parse(__payload)
	if !_ok {
		slog.Error("pay: bad coin payload", "payload", __payload, "charge", __charge)
		return nil
	}
	_plan, _until := _p.ID, time.Now().Add(pay.Period*time.Second)
	_old, _ := __d.Subs.Get(__ctx, _u)
	if _old.Live() {
		_plan, _until = pay.Merge(_old.Plan, _old.Until, _p.ID, pay.Period*time.Second, time.Now())
	}
	_amt, _cur := _p.Rub, "RUB"
	if strings.HasPrefix(__charge, "xr:") {
		_amt, _cur = cents(_p.Usd), "USDT/100"
	}
	_new, _err := __d.Subs.Pay(__ctx, _u, _plan, _amt, _cur, _until, __charge)
	if _err != nil || !_new {
		return _err
	}
	if _old.Live() && pay.Star(_old.Charge) {
		if _err = __d.Tg.Cancel(__ctx, _u, _old.Charge); _err != nil {
			slog.Error("pay: cancel old sub", "user", _u, "err", _err)
		}
	}
	_n := _p.Name
	if _q, _ok := pay.Get(_plan); _ok {
		_n = _q.Name
	}
	slog.Info("pay", "user", _u, "plan", _plan, "via", __charge, "until", _until)
	return __d.Tg.Send(__ctx, _u, text.PaidCoin(_n, _until), kb.One(kb.Web(text.SetBtn, __d.Url("/")).With(text.Gear)))
}

func Remind(__ctx context.Context, __d *deps.D) {
	for {
		_ds, _err := __d.Subs.Due(__ctx, 72*time.Hour)
		if _err != nil {
			slog.Warn("remind", "err", _err)
		}
		for _, _x := range _ds {
			_p, _ := pay.Get(_x.Plan)
			_ = __d.Tg.Send(__ctx, _x.User, text.Renew(_p.Name, _x.Until), kb.One(kb.Web(text.RenewBtn, __d.Url("/f/plans?to="+_x.Plan)).With(text.Spark)))
		}
		select {
		case <-__ctx.Done():
			return
		case <-time.After(time.Hour):
		}
	}
}

func cents(__s string) int {
	_a, _b, _ := strings.Cut(__s, ".")
	_n, _ := strconv.Atoi(_a + (_b + "00")[:2])
	return _n
}

func Sweep(__ctx context.Context, __d *deps.D) {
	for {
		if __d.Cfg.Cb != "" {
			if _l, _err := cb.Paid(__ctx, __d.Cfg.Cb); _err != nil {
				slog.Warn("sweep cb", "err", _err)
			} else {
				for _, _i := range _l {
					if _i.Status == "paid" && _i.Payload != "" {
						_ = Coin(__ctx, __d, _i.Payload, "cb:"+strconv.FormatInt(_i.ID, 10))
					}
				}
			}
		}
		if __d.Cfg.Xr != "" {
			if _l, _err := xr.Paid(__ctx, __d.Cfg.Xr); _err != nil {
				slog.Warn("sweep xr", "err", _err)
			} else {
				for _, _i := range _l {
					if _i.Status == "paid" && _i.ClientID != "" {
						_ = Coin(__ctx, __d, xr.Client(_i.ClientID), "xr:"+_i.ID)
					}
				}
			}
		}
		select {
		case <-__ctx.Done():
			return
		case <-time.After(2 * time.Minute):
		}
	}
}
