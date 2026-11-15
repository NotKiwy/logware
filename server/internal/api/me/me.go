package me

import (
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/pay"
	"logware/internal/store/users"
)

type Conn struct {
	Active bool `json:"active"`
	Reply  bool `json:"reply"`
	Read   bool `json:"read"`
	DelOwn bool `json:"delOwn"`
	DelAll bool `json:"delAll"`
	Name   bool `json:"name"`
}

type Stats struct {
	Caught int `json:"caught"`
	Today  int `json:"today"`
	Kept   int `json:"kept"`
}

type Me struct {
	Plan  string      `json:"plan"`
	Until int64       `json:"until"`
	Stats Stats       `json:"stats"`
	Reg   bool        `json:"reg"`
	Mode  string      `json:"mode,omitempty"`
	Vault []byte      `json:"vault,omitempty"`
	Prefs users.Prefs `json:"prefs"`
	Conn  *Conn       `json:"conn"`
	Bot   string      `json:"bot"`
	Dev   bool        `json:"dev"`
	Via   []string    `json:"via"`
	Stars bool        `json:"stars"`
	Sup   string      `json:"support"`
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := auth.ID(r)
		_u, _err := __d.Users.Get(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_c, _err := __d.Conns.ByUser(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_me := Me{Plan: "free", Prefs: users.Def(), Bot: __d.Nick, Sup: __d.Sup}
		if _u.Ready() {
			_me.Reg, _me.Mode, _me.Prefs = true, _u.Mode, _u.Prefs
			if _u.Mode == "pass" {
				_me.Vault = _u.Vault
			}
		}
		if _me.Reg {
			_me.Stats.Caught, _me.Stats.Today, _ = __d.Evs.Count(r.Context(), _id)
			_me.Stats.Kept, _ = __d.Msgs.Count(r.Context(), _id)
		}
		if _s, _ := __d.Subs.Get(r.Context(), _id); _s.Live() {
			_me.Plan, _me.Until, _me.Stars = _s.Plan, _s.Until.Unix(), pay.Star(_s.Charge)
		}
		_me.Via = []string{"stars"}
		if __d.Cfg.Cb != "" {
			_me.Via = append(_me.Via, "cb")
		}
		if __d.Cfg.Xr != "" && __d.Cfg.XrWh != "" {
			_me.Via = append(_me.Via, "xr")
		}
		_me.Dev = _me.Plan == "dev"
		if _c != nil {
			_me.Conn = &Conn{Active: _c.Active, Reply: _c.Rights.Reply, Read: _c.Rights.Read, DelOwn: _c.Rights.DelOwn, DelAll: _c.Rights.DelAll, Name: _c.Rights.Name}
		}
		out.JSON(w, _me)
	}
}
