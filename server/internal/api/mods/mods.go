package mods

import (
	"context"
	"net/http"
	"strconv"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/mods/mod"
	"logware/internal/mods/reg"
)

type View struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Desc   string     `json:"desc"`
	Long   string     `json:"long"`
	Cat    string     `json:"cat"`
	Ver    string     `json:"ver"`
	Check  bool       `json:"check"`
	Perms  []mod.Perm `json:"perms"`
	Icon   string     `json:"icon"`
	Color  string     `json:"color"`
	Author string     `json:"author"`
	Pic    string     `json:"pic"`
	Soon   bool       `json:"soon"`
	Hidden bool       `json:"hidden"`
	On     bool       `json:"on"`
	Ready  bool       `json:"ready"`
	Cfg    mod.Cfg    `json:"cfg"`
	Opts   []mod.Opt  `json:"opts"`
}

type In struct {
	On  bool              `json:"on"`
	Cfg map[string]string `json:"cfg"`
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_st, _err := __d.Mods.List(r.Context(), auth.ID(r))
		_pro := __d.Pro(r.Context(), auth.ID(r))
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_lv := Live(r.Context(), __d)
		_out := make([]View, 0, len(reg.All()))
		for _, _m := range reg.All() {
			_s := _st[_m.ID]
			_opts := _m.Opts
			if _opts == nil {
				_opts = []mod.Opt{}
			}
			_perms := _m.Perms
			if _perms == nil {
				_perms = []mod.Perm{}
			}
			_v := View{
				ID: _m.ID, Name: _m.Name, Desc: _m.Desc, Long: _m.Long, Cat: _m.Cat, Ver: _m.Ver, Check: _m.Check, Perms: _perms, Icon: _m.Icon, Color: _m.Color, Author: _m.Author, Soon: _m.Soon, Hidden: _m.Hidden,
				On: _s.On && _m.Can(__d) && _pro, Ready: _m.Can(__d), Cfg: _m.Fill(_s.Cfg), Opts: _opts,
			}
			if _l, _ok := _lv[_m.ID]; _ok {
				_v.Name, _v.Desc, _v.Ver, _v.Author, _v.Pic = _l.Name, _l.Desc, _l.Ver, _l.Author, _l.Pic
				if _l.Desc != "" {
					_v.Long = _l.Desc
				}
			}
			_out = append(_out, _v)
		}
		out.JSON(w, _out)
	}
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_m := reg.Get(r.PathValue("id"))
		if _m == nil {
			out.Fail(w, http.StatusNotFound, "no such module")
			return
		}
		var _in In
		if !out.Read(w, r, &_in) {
			return
		}
		_u, _err := __d.Users.Get(r.Context(), auth.ID(r))
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if !_u.Ready() {
			out.Fail(w, http.StatusConflict, "not registered")
			return
		}
		_cfg := _m.Fill(_in.Cfg)
		_on := _in.On && _m.Can(__d)
		if _on && !__d.Pro(r.Context(), _u.ID) {
			out.Fail(w, http.StatusPaymentRequired, "need pro")
			return
		}
		if _m.Set != nil {
			if _cfg, _err = _m.Set(r.Context(), __d, _u.ID, _on, _cfg); _err != nil {
				out.Fail(w, http.StatusConflict, _err.Error())
				return
			}
			_cfg = _m.Fill(_cfg)
		}
		if _err = __d.Mods.Set(r.Context(), _u.ID, _m.ID, _on, _cfg); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]any{"on": _on, "cfg": _cfg})
	}
}

type Info struct {
	ID     int64
	Name   string
	Desc   string
	Ver    string
	Author string
	Pic    string
}

func PicURL(__id int64, __v int) string {
	if __v == 0 {
		return ""
	}
	return "/pic/" + strconv.FormatInt(__id, 10) + "?v=" + strconv.Itoa(__v)
}

func Live(__ctx context.Context, __d *deps.D) map[string]Info {
	_out := map[string]Info{}
	_ps, _err := __d.Devs.Live(__ctx)
	if _err != nil {
		return _out
	}
	_ids := make([]int64, 0, len(_ps))
	for _, _p := range _ps {
		_ids = append(_ids, _p.User)
	}
	_nicks, _ := __d.Devs.Nicks(__ctx, _ids)
	for _, _p := range _ps {
		_out[_p.Slug] = Info{ID: _p.ID, Name: _p.Name, Desc: _p.About, Ver: _p.Ver, Author: _nicks[_p.User], Pic: PicURL(_p.ID, _p.Pv)}
	}
	return _out
}

func Pic(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		_b, _err := __d.Devs.Pic(r.Context(), _id)
		if _err != nil || _b == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", http.DetectContentType(_b))
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(_b)
	}
}
