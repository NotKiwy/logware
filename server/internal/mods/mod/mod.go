package mod

import (
	"context"

	"logware/internal/bot/own"
	"logware/internal/deps"
	"logware/internal/tg/types"
	"logware/sdk"
)

type (
	Cfg    = sdk.Cfg
	Choice = sdk.Choice
	Opt    = sdk.Opt
	Perm   = sdk.Perm
)

type Mod struct {
	ID     string
	Name   string
	Desc   string
	Long   string
	Cat    string
	Ver    string
	Check  bool
	Perms  []Perm
	Icon   string
	Color  string
	Author string
	Soon   bool
	Hidden bool
	Opts   []Opt
	Ready  func(*deps.D) bool
	Msg    func(context.Context, *deps.D, *own.Own, *types.Message, Cfg) error
	Set    func(context.Context, *deps.D, int64, bool, Cfg) (Cfg, error)
}

func (m *Mod) Fill(__c map[string]string) Cfg {
	_out := Cfg{}
	for _, _o := range m.Opts {
		_out[_o.Key] = _o.Def
		if _v, _ok := __c[_o.Key]; _ok && _o.Has(_v) {
			_out[_o.Key] = _v
		}
	}
	return _out
}

func (m *Mod) Can(__d *deps.D) bool {
	return !m.Soon && m.Ready != nil && m.Ready(__d)
}
