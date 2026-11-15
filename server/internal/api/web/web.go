package web

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var __routes__ = regexp.MustCompile(`^/(|reg|mods|set|dev|dev/new|dev/p/\d+|dev/p/\d+/(update|edit)|f/[a-z]+|p/[a-z0-9-]+|e/[a-f0-9]+)$`)

func Static(__dir string) http.Handler {
	_fs := http.FileServer(http.Dir(__dir))
	_index := filepath.Join(__dir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}
		_p := filepath.Join(__dir, filepath.Clean("/"+r.URL.Path))
		if _st, _err := os.Stat(_p); _err == nil && !_st.IsDir() {
			_fs.ServeHTTP(w, r)
			return
		}
		if !__routes__.MatchString(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, _index)
	})
}
