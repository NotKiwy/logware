package route

import (
	"context"
	"strings"

	"logware/internal/bot/once"
	"logware/internal/bot/own"
	"logware/internal/deps"
	"logware/internal/tg/types"
)

type Fn func(context.Context, *deps.D, *own.Own, *types.Message, []string) error

var __cmds__ = map[string]Fn{
	"s": once.Cmd,
}

func IsCmd(__text string) bool {
	return len(__text) > 1 && __text[0] == '.' && __text[1] != '.'
}

func Run(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message) error {
	_f := strings.Fields(strings.TrimPrefix(__m.Text, "."))
	if len(_f) == 0 {
		return nil
	}
	_fn, _ok := __cmds__[strings.ToLower(_f[0])]
	if !_ok {
		return nil
	}
	return _fn(__ctx, __d, __o, __m, _f[1:])
}
