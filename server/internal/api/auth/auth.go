package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TTL = 24 * time.Hour

var ErrBad = errors.New("bad init data")

type _key struct{}

type Who struct {
	ID   int64  `json:"id"`
	Name string `json:"first_name"`
	Last string `json:"last_name"`
	Nick string `json:"username"`
	Prem bool   `json:"is_premium"`
}

func (w Who) Full() string {
	return strings.TrimSpace(w.Name + " " + w.Last)
}

func Check(__tok, __raw string, __now time.Time) (int64, error) {
	_w, _err := Parse(__tok, __raw, __now)
	return _w.ID, _err
}

func Parse(__tok, __raw string, __now time.Time) (Who, error) {
	var _w Who
	_v, _err := url.ParseQuery(__raw)
	if _err != nil {
		return _w, ErrBad
	}
	_hash := _v.Get("hash")
	if _hash == "" {
		return _w, ErrBad
	}
	_keys := make([]string, 0, len(_v))
	for _k := range _v {
		if _k != "hash" {
			_keys = append(_keys, _k)
		}
	}
	sort.Strings(_keys)
	_parts := make([]string, len(_keys))
	for _i, _k := range _keys {
		_parts[_i] = _k + "=" + _v.Get(_k)
	}
	if !hmac.Equal([]byte(Sign(__tok, strings.Join(_parts, "\n"))), []byte(_hash)) {
		return _w, ErrBad
	}
	_at, _err := strconv.ParseInt(_v.Get("auth_date"), 10, 64)
	if _err != nil || __now.Sub(time.Unix(_at, 0)) > TTL {
		return _w, ErrBad
	}
	if json.Unmarshal([]byte(_v.Get("user")), &_w) != nil || _w.ID == 0 {
		return Who{}, ErrBad
	}
	return _w, nil
}

func Sign(__tok, __data string) string {
	_sec := hmac.New(sha256.New, []byte("WebAppData"))
	_sec.Write([]byte(__tok))
	_m := hmac.New(sha256.New, _sec.Sum(nil))
	_m.Write([]byte(__data))
	return hex.EncodeToString(_m.Sum(nil))
}

func Mw(__tok string, __ok func(context.Context, Who) bool) func(http.Handler) http.Handler {
	return func(__next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_w, _err := Parse(__tok, r.Header.Get("X-Init"), time.Now())
			if _err != nil {
				http.NotFound(w, r)
				return
			}
			if __ok != nil && !__ok(r.Context(), _w) {
				http.Error(w, "banned", http.StatusForbidden)
				return
			}
			__next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), _key{}, _w.ID)))
		})
	}
}

func ID(r *http.Request) int64 {
	_id, _ := r.Context().Value(_key{}).(int64)
	return _id
}
