package limit

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type _Bucket struct {
	left float64
	at   time.Time
}

type Limit struct {
	mu    sync.Mutex
	rate  float64
	burst float64
	all   map[string]*_Bucket
	swept time.Time
}

func New(__rate, __burst float64) *Limit {
	return &Limit{rate: __rate, burst: __burst, all: map[string]*_Bucket{}}
}

func (l *Limit) Allow(__key string, __now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if __now.Sub(l.swept) > time.Minute {
		for _k, _b := range l.all {
			if __now.Sub(_b.at).Seconds()*l.rate > l.burst {
				delete(l.all, _k)
			}
		}
		l.swept = __now
	}
	_b, _ok := l.all[__key]
	if !_ok {
		_b = &_Bucket{left: l.burst, at: __now}
		l.all[__key] = _b
	}
	_b.left += __now.Sub(_b.at).Seconds() * l.rate
	if _b.left > l.burst {
		_b.left = l.burst
	}
	_b.at = __now
	if _b.left < 1 {
		return false
	}
	_b.left--
	return true
}

func IP(r *http.Request) string {
	if _f := r.Header.Get("X-Forwarded-For"); _f != "" {
		_ip, _, _ := strings.Cut(_f, ",")
		return strings.TrimSpace(_ip)
	}
	_h, _, _err := net.SplitHostPort(r.RemoteAddr)
	if _err != nil {
		return r.RemoteAddr
	}
	return _h
}

func Mw(__l *Limit, __key func(*http.Request) string) func(http.Handler) http.Handler {
	return func(__next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !__l.Allow(__key(r), time.Now()) {
				w.Header().Set("Retry-After", strconv.Itoa(int(1/__l.rate)+1))
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			__next.ServeHTTP(w, r)
		})
	}
}
