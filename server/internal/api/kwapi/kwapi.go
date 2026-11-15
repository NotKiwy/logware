package kwapi

import (
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/store/kws"
)

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_l, _err := __d.Kws.List(r.Context(), auth.ID(r))
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, _l)
	}
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _k kws.Kw
		if !out.Read(w, r, &_k) {
			return
		}
		if len(_k.Hash) != 32 || len(_k.Box) == 0 || len(_k.Box) > 512 {
			out.Fail(w, http.StatusBadRequest, "bad word")
			return
		}
		_ok, _err := __d.Kws.Add(r.Context(), auth.ID(r), _k)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if !_ok {
			out.Fail(w, http.StatusConflict, "limit")
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}

func Del(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _k kws.Kw
		if !out.Read(w, r, &_k) {
			return
		}
		if _err := __d.Kws.Del(r.Context(), auth.ID(r), _k.Hash); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
