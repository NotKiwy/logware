package clock

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"logware/internal/deps"
	"logware/internal/mods/mod"
	"logware/internal/tg/client"
)

var (
	__mu   sync.Mutex
	__last = map[int64]string{}
)

func forget(__user int64) {
	__mu.Lock()
	delete(__last, __user)
	__mu.Unlock()
}

func Run(__ctx context.Context, __d *deps.D) {
	for {
		_wait := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute + time.Second))
		select {
		case <-__ctx.Done():
			return
		case <-time.After(_wait):
		}
		tick(__ctx, __d)
	}
}

func tick(__ctx context.Context, __d *deps.D) {
	_rows, _err := __d.Mods.Active(__ctx, ID)
	if _err != nil {
		slog.Error("clock", "err", _err)
		return
	}
	for _, _r := range _rows {
		if !__d.Pro(__ctx, _r.User) {
			if _c, _err := set(__ctx, __d, _r.User, false, Mod.Fill(_r.Cfg)); _err == nil {
				_ = __d.Mods.Set(__ctx, _r.User, ID, false, _c)
			}
			continue
		}
		if _err = one(__ctx, __d, _r.User, Mod.Fill(_r.Cfg)); _err != nil {
			slog.Warn("clock", "user", _r.User, "err", _err)
			var _te *client.Err
			if errors.As(_err, &_te) && _te.Code == 429 {
				return
			}
		}
	}
}

func one(__ctx context.Context, __d *deps.D, __user int64, __c mod.Cfg) error {
	if __c["first"] == "" {
		return nil
	}
	_c, _err := __d.Conns.ByUser(__ctx, __user)
	if _err != nil || _c == nil || !_c.Active || !_c.Rights.Name {
		return _err
	}
	_f, _l := Name(__c, now(__c))
	_k := Key(_f, _l)
	__mu.Lock()
	_same := __last[__user] == _k
	__mu.Unlock()
	if _same {
		return nil
	}
	if _err = __d.Tg.BizName(__ctx, _c.ID, _f, _l); _err != nil {
		return _err
	}
	__mu.Lock()
	__last[__user] = _k
	__mu.Unlock()
	return nil
}
