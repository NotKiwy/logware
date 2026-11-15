package poll

import (
	"context"
	"log/slog"
	"time"

	"logware/internal/bot/upd"
	"logware/internal/deps"
)

func Run(__ctx context.Context, __d *deps.D) error {
	var _off int64
	for {
		_us, _err := __d.Tg.Updates(__ctx, _off)
		if __ctx.Err() != nil {
			return nil
		}
		if _err != nil {
			slog.Error("poll", "err", _err)
			select {
			case <-__ctx.Done():
				return nil
			case <-time.After(3 * time.Second):
			}
			continue
		}
		for _, _u := range _us {
			_off = _u.ID + 1
			if _err = upd.Run(__ctx, __d, &_u); _err != nil {
				slog.Error("update", "id", _u.ID, "err", _err)
			}
		}
	}
}
