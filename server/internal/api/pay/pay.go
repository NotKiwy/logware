package pay

import (
	"net/http"

	"strconv"
	"time"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/coin/cb"
	"logware/internal/coin/xr"
	"logware/internal/deps"
	"logware/internal/pay"
)

type In struct {
	Plan string `json:"plan"`
	Via  string `json:"via"`
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _in In
		if !out.Read(w, r, &_in) {
			return
		}
		_p, _ok := pay.Get(_in.Plan)
		if !_ok {
			out.Fail(w, http.StatusBadRequest, "no such plan")
			return
		}
		_id := auth.ID(r)
		_u, _err := __d.Users.Get(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if !_u.Ready() {
			out.Fail(w, http.StatusConflict, "not registered")
			return
		}
		_s, _err := __d.Subs.Get(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_back := "https://t.me/" + __d.Nick
		switch _in.Via {
		case "cb", "xr":
			if _s.Live() && pay.Rank(_s.Plan) > pay.Rank(_p.ID) {
				out.Fail(w, http.StatusConflict, "already paid")
				return
			}
			if _s.Live() && pay.Star(_s.Charge) {
				out.Fail(w, http.StatusConflict, "has stars")
				return
			}
			var _url string
			if _in.Via == "cb" && __d.Cfg.Cb != "" {
				_i, _err := cb.Create(r.Context(), __d.Cfg.Cb, _p.Rub, _p.Name, pay.Tagged(_p.ID, _id, __d.Cfg.Tag), _back)
				if _err != nil {
					out.Oops(w, _err)
					return
				}
				_url = _i.Link()
			} else if _in.Via == "xr" && __d.Cfg.Xr != "" {
				_i, _err := xr.Create(r.Context(), __d.Cfg.Xr, _p.Usd, _p.Name, pay.Tagged(_p.ID, _id, __d.Cfg.Tag)+":"+strconv.FormatInt(time.Now().UnixNano(), 36), _id, _back)
				if _err != nil {
					out.Oops(w, _err)
					return
				}
				_url = _i.Links.Bot
			} else {
				out.Fail(w, http.StatusBadRequest, "no method")
				return
			}
			out.JSON(w, map[string]string{"url": _url})
			return
		}
		if _s.Live() && pay.Rank(_s.Plan) >= pay.Rank(_p.ID) {
			out.Fail(w, http.StatusConflict, "already paid")
			return
		}
		_url, _err := __d.Tg.Invoice(r.Context(), _p.Name, _p.Name, pay.Payload(_p.ID, _id), _p.Price, pay.Period)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]string{"url": _url})
	}
}
