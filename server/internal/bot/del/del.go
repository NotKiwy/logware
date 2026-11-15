package del

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"logware/internal/bot/notify"
	"logware/internal/bot/own"
	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/types"
)

const (
	__wait__ = 3 * time.Second
	__cap__  = 30 * time.Second
)

type _Key struct {
	conn string
	chat int64
}

type _Box struct {
	d     *deps.D
	o     *own.Own
	name  string
	refs  []int64
	seen  map[int64]bool
	first time.Time
	t     *time.Timer
}

var (
	__mu    sync.Mutex
	__boxes = map[_Key]*_Box{}
)

func Run(__ctx context.Context, __d *deps.D, __x *types.Deleted) error {
	_o, _err := own.Get(__ctx, __d, __x.Conn)
	if _err != nil || !_o.Ready() || !_o.User.Prefs.Del {
		return _err
	}
	if _skip, _err := __d.Ign.Has(__ctx, _o.User.ID, __x.Chat.ID); _err != nil || _skip {
		return _err
	}
	_rows, _err := __d.Msgs.Find(__ctx, __x.Conn, __x.Chat.ID, __x.IDs)
	if _err != nil {
		return _err
	}
	var _refs []int64
	var _msgs []int64
	for _, _r := range _rows {
		if _r.Out && !_o.User.Prefs.DelMe {
			continue
		}
		_refs = append(_refs, _r.ID)
		_msgs = append(_msgs, _r.Msg)
	}
	if len(_refs) == 0 {
		return nil
	}
	put(_Key{__x.Conn, __x.Chat.ID}, __d, _o, text.Name(__x.Chat), _refs, _msgs)
	return nil
}

func put(__k _Key, __d *deps.D, __o *own.Own, __name string, __refs, __msgs []int64) {
	__mu.Lock()
	defer __mu.Unlock()
	_b := __boxes[__k]
	if _b == nil {
		_b = &_Box{d: __d, name: __name, seen: map[int64]bool{}, first: time.Now()}
		_b.t = time.AfterFunc(__wait__, func() { flush(__k) })
		__boxes[__k] = _b
	} else if time.Since(_b.first) < __cap__ {
		_b.t.Reset(__wait__)
	}
	_b.o = __o
	_b.refs = append(_b.refs, __refs...)
	for _, _m := range __msgs {
		_b.seen[_m] = true
	}
}

func flush(__k _Key) {
	__mu.Lock()
	_b := __boxes[__k]
	delete(__boxes, __k)
	__mu.Unlock()
	if _b == nil {
		return
	}
	_ctx, _done := context.WithTimeout(context.Background(), 30*time.Second)
	defer _done()
	_n := len(_b.seen)
	_p := _b.o.User.Prefs
	var _err error
	if _n < 50 && _p.Panic && _n >= _p.PanicN {
		_err = notify.Loud(_ctx, _b.d, _b.o, "del", _b.refs, text.Panic(_b.name, _n))
	} else {
		_err = notify.Send(_ctx, _b.d, _b.o, "del", _b.refs, text.Del(_b.name, _n))
	}
	if _err != nil {
		slog.Error("del notify", "err", _err)
	}
}
