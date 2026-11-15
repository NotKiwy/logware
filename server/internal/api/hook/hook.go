package hook

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"logware/internal/bot/pay"
	"logware/internal/coin/cb"
	"logware/internal/coin/xr"
	"logware/internal/deps"
)

const __max__ = 64 << 10

func body(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	_b, _err := io.ReadAll(io.LimitReader(r.Body, __max__+1))
	if _err != nil || len(_b) > __max__ {
		http.Error(w, "bad body", http.StatusBadRequest)
		return nil, false
	}
	return _b, true
}

func Cb(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_b, _ok := body(w, r)
		if !_ok {
			return
		}
		slog.Info("hook cb", "len", len(_b))
		if !cb.Verify(__d.Cfg.Cb, _b, r.Header.Get("crypto-pay-api-signature")) {
			slog.Warn("hook cb: bad signature")
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		var _u cb.Update
		if json.Unmarshal(_b, &_u) != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if _u.Type == "invoice_paid" && _u.Inv.Status == "paid" {
			if _err := pay.Coin(r.Context(), __d, _u.Inv.Payload, "cb:"+itoa(_u.Inv.ID)); _err != nil {
				slog.Error("cb hook", "err", _err)
				http.Error(w, "retry", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func Xr(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_b, _ok := body(w, r)
		if !_ok {
			return
		}
		slog.Info("hook xr", "len", len(_b))
		if !xr.Verify(__d.Cfg.XrWh, _b, r.Header.Get("Signature"), r.Header.Get("Signature-Version"), r.Header.Get("Signature-Timestamp"), time.Now()) {
			slog.Warn("hook xr: bad signature")
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		var _e xr.Event
		if json.Unmarshal(_b, &_e) != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if _e.Type == "invoice" && _e.Data.Event == "invoice_status_changed" && _e.Data.Inv.Status == "paid" {
			if _err := pay.Coin(r.Context(), __d, xr.Client(_e.Data.Inv.ClientID), "xr:"+_e.Data.Inv.ID); _err != nil {
				slog.Error("xr hook", "err", _err)
				http.Error(w, "retry", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func itoa(__n int64) string {
	return strconv.FormatInt(__n, 10)
}

func Any(__d *deps.D) http.HandlerFunc {
	_cb, _xr := Cb(__d), Xr(__d)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("crypto-pay-api-signature") != "" {
			_cb(w, r)
			return
		}
		_xr(w, r)
	}
}
