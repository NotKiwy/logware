package clean

import (
	"context"
	"log/slog"
	"time"

	"logware/internal/deps"
)

func Run(__ctx context.Context, __d *deps.D) {
	_t := time.NewTicker(time.Hour)
	defer _t.Stop()
	for {
		once(__ctx, __d)
		select {
		case <-__ctx.Done():
			return
		case <-_t.C:
		}
	}
}

func once(__ctx context.Context, __d *deps.D) {
	_m, _err := __d.Msgs.Clean(__ctx)
	if _err != nil {
		slog.Error("clean msgs", "err", _err)
		return
	}
	_e, _err := __d.Evs.Clean(__ctx)
	if _err != nil {
		slog.Error("clean events", "err", _err)
		return
	}
	if _m+_e > 0 {
		slog.Info("cleaned", "msgs", _m, "events", _e)
	}
}
