package reg

import (
	"net/http"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/crypto/seal"
	"logware/internal/deps"
)

var __modes__ = map[string]bool{"cloud": true, "secure": true, "pass": true}

type In struct {
	Mode  string `json:"mode"`
	Pub   []byte `json:"pub"`
	Vault []byte `json:"vault"`
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _in In
		if !out.Read(w, r, &_in) {
			return
		}
		if !__modes__[_in.Mode] || len(_in.Pub) != seal.KeyLen {
			out.Fail(w, http.StatusBadRequest, "bad mode or key")
			return
		}
		if (_in.Mode == "pass") != (len(_in.Vault) > 0) || len(_in.Vault) > 1024 {
			out.Fail(w, http.StatusBadRequest, "bad vault")
			return
		}
		_id, _ctx := auth.ID(r), r.Context()
		if _err := __d.Msgs.Drop(_ctx, _id); _err != nil {
			out.Oops(w, _err)
			return
		}
		if _err := __d.Evs.Drop(_ctx, _id); _err != nil {
			out.Oops(w, _err)
			return
		}
		if _err := __d.Users.Reg(_ctx, _id, _in.Mode, _in.Pub, _in.Vault); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
