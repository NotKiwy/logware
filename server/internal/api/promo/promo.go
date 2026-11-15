package promo

import (
	"errors"
	"log/slog"
	"net/http"
	"unicode/utf8"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/pay"
	"logware/internal/store/promos"
)

type In struct {
	Code string `json:"code"`
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _in In
		if !out.Read(w, r, &_in) {
			return
		}
		_c := promos.Norm(_in.Code)
		if _c == "" || utf8.RuneCountInString(_c) > 32 {
			out.Fail(w, http.StatusBadRequest, "none")
			return
		}
		_id := auth.ID(r)
		_plan, _until, _err := __d.Promos.Use(r.Context(), _c, _id, pay.Merge)
		switch {
		case errors.Is(_err, promos.ErrNone):
			out.Fail(w, http.StatusBadRequest, "none")
			return
		case errors.Is(_err, promos.ErrUsed):
			out.Fail(w, http.StatusBadRequest, "used")
			return
		case errors.Is(_err, promos.ErrOver):
			out.Fail(w, http.StatusBadRequest, "over")
			return
		case _err != nil:
			out.Oops(w, _err)
			return
		}
		slog.Info("promo", "user", _id, "code", _c, "plan", _plan, "until", _until)
		out.JSON(w, map[string]any{"plan": _plan, "until": _until.Unix()})
	}
}
