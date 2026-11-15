package chats

import (
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
)

type _Row struct {
	ID   int64  `json:"id"`
	Body []byte `json:"body"`
	N    int    `json:"n"`
	Off  bool   `json:"off"`
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := auth.ID(r)
		_cs, _err := __d.Msgs.Chats(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_off, _err := __d.Ign.List(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_out := make([]_Row, 0, len(_cs))
		for _, _c := range _cs {
			_out = append(_out, _Row{ID: _c.ID, Body: _c.Body, N: _c.N, Off: _off[_c.ID]})
		}
		out.JSON(w, _out)
	}
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _in struct {
			ID  int64 `json:"id"`
			Off bool  `json:"off"`
		}
		if !out.Read(w, r, &_in) {
			return
		}
		if _err := __d.Ign.Set(r.Context(), auth.ID(r), _in.ID, _in.Off); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
