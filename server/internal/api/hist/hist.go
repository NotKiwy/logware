package hist

import (
	"net/http"
	"strconv"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
)

func num(__r *http.Request, __k string, __def, __max int) int {
	_n, _err := strconv.Atoi(__r.URL.Query().Get(__k))
	if _err != nil || _n < 0 {
		return __def
	}
	return min(_n, __max)
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_lim := num(r, "n", 40, 200)
		_rows, _err := __d.Evs.List(r.Context(), auth.ID(r), num(r, "from", 0, 1<<30), _lim+1)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_more := len(_rows) > _lim
		if _more {
			_rows = _rows[:_lim]
		}
		out.JSON(w, map[string]any{"items": _rows, "more": _more})
	}
}

func Clear(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id, _ctx := auth.ID(r), r.Context()
		if _err := __d.Msgs.Drop(_ctx, _id); _err != nil {
			out.Oops(w, _err)
			return
		}
		if _err := __d.Evs.Drop(_ctx, _id); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
