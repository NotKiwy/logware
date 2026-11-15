package prefs

import (
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/store/users"
)

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_p := users.Def()
		if !out.Read(w, r, &_p) {
			return
		}
		_p = _p.Fix()
		_u, _err := __d.Users.Get(r.Context(), auth.ID(r))
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if !_u.Ready() {
			out.Fail(w, http.StatusConflict, "not registered")
			return
		}
		if _err = __d.Users.SetPrefs(r.Context(), _u.ID, _p); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, _p)
	}
}
