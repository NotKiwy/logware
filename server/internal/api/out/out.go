package out

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(w http.ResponseWriter, __v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(__v)
}

func Fail(w http.ResponseWriter, __code int, __msg string) {
	http.Error(w, __msg, __code)
}

func Oops(w http.ResponseWriter, __err error) {
	slog.Error("api", "err", __err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func Read(w http.ResponseWriter, r *http.Request, __v any) bool {
	if _err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(__v); _err != nil {
		Fail(w, http.StatusBadRequest, "bad json")
		return false
	}
	return true
}
