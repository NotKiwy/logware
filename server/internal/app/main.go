package app

import (
	_ "time/tzdata"

	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Main() {
	_ctx, _stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer _stop()
	if _err := Run(_ctx); _err != nil {
		slog.Error("logware", "err", _err)
		os.Exit(1)
	}
}
