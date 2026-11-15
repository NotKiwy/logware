package reg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"logware/internal/bot/own"
	"logware/internal/deps"
	"logware/internal/mods/clock"
	"logware/internal/mods/mod"
	"logware/internal/tg/types"
	"logware/sdk"
)

var __all__ = []*mod.Mod{clock.Mod}

func Add(__p ...*sdk.Plugin) {
	__all__ = append(__all__, wrap(__p)...)
}

var ErrRight = errors.New("plugin has no such right")

type _Bot struct {
	d *deps.D
	o *own.Own
	m *types.Message
	p *sdk.Plugin
}

func (b *_Bot) Edit(__ctx context.Context, __t string) error {
	if !b.p.Can(sdk.Reply) {
		return ErrRight
	}
	return b.d.Tg.EditBiz(__ctx, b.o.Conn.ID, b.m.Chat.ID, b.m.ID, __t)
}

func (b *_Bot) Reply(__ctx context.Context, __t string) error {
	if !b.p.Can(sdk.Reply) {
		return ErrRight
	}
	return b.d.Tg.SendBiz(__ctx, b.o.Conn.ID, b.m.Chat.ID, __t, b.m.ID)
}

func (b *_Bot) Note(__ctx context.Context, __t string) error {
	return b.d.Tg.Send(__ctx, b.o.Conn.Chat, __t, nil)
}

func wrap(__ps []*sdk.Plugin) []*mod.Mod {
	_out := make([]*mod.Mod, 0, len(__ps))
	for _, _p := range __ps {
		_out = append(_out, &mod.Mod{
			ID: _p.ID, Name: _p.Name, Desc: _p.Desc, Long: _p.Long, Cat: _p.Cat, Ver: _p.Ver, Check: _p.Check,
			Perms: _p.Perms, Icon: _p.Icon, Color: _p.Color, Author: _p.Author, Opts: _p.Opts,
			Ready: func(*deps.D) bool { return true },
			Msg: func(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message, __c mod.Cfg) (_err error) {
				if _p.On == nil {
					return nil
				}
				defer func() {
					if _r := recover(); _r != nil {
						_err = fmt.Errorf("plugin %s panicked: %v", _p.ID, _r)
					}
				}()
				__ctx, _stop := context.WithTimeout(__ctx, 10*time.Second)
				defer _stop()
				_t := __m.Text
				if _t == "" {
					_t = __m.Caption
				}
				return _p.On(__ctx, &_Bot{d: __d, o: __o, m: __m, p: _p}, sdk.Msg{Text: _t, Kind: __m.Kind(), Out: __m.Out(__o.Conn.User)}, __c)
			},
		})
	}
	return _out
}

func All() []*mod.Mod {
	return __all__
}

func Get(__id string) *mod.Mod {
	for _, _m := range __all__ {
		if _m.ID == __id {
			return _m
		}
	}
	return nil
}

func Msg(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message) {
	_st, _err := __d.Mods.List(__ctx, __o.User.ID)
	if _err != nil {
		slog.Error("mods", "err", _err)
		return
	}
	_pro := false
	for _, _s := range _st {
		_pro = _pro || _s.On
	}
	if !_pro || !__d.Pro(__ctx, __o.User.ID) {
		return
	}
	for _, _m := range __all__ {
		_s, _ok := _st[_m.ID]
		if !_ok || !_s.On || _m.Msg == nil || !_m.Can(__d) {
			continue
		}
		if _err = _m.Msg(__ctx, __d, __o, __m, _m.Fill(_s.Cfg)); _err != nil {
			slog.Warn("mod", "id", _m.ID, "err", _err)
		}
	}
}
