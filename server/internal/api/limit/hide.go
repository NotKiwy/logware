package limit

import "net/http"

type _Hide struct {
	http.ResponseWriter
	gone bool
}

func (h *_Hide) WriteHeader(__code int) {
	if __code == http.StatusMethodNotAllowed {
		h.gone = true
		h.Header().Del("Allow")
		h.Header().Set("Content-Type", "text/plain; charset=utf-8")
		h.ResponseWriter.WriteHeader(http.StatusNotFound)
		_, _ = h.ResponseWriter.Write([]byte("404 page not found\n"))
		return
	}
	h.ResponseWriter.WriteHeader(__code)
}

func (h *_Hide) Write(__b []byte) (int, error) {
	if h.gone {
		return len(__b), nil
	}
	return h.ResponseWriter.Write(__b)
}

func Hide(__next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		__next.ServeHTTP(&_Hide{ResponseWriter: w}, r)
	})
}
