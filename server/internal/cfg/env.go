package cfg

import (
	"bufio"
	"os"
	"strings"
)

func Env(__paths ...string) {
	for _, _p := range __paths {
		_f, _err := os.Open(_p)
		if _err != nil {
			continue
		}
		_sc := bufio.NewScanner(_f)
		for _sc.Scan() {
			_line := strings.TrimSpace(_sc.Text())
			if _line == "" || strings.HasPrefix(_line, "#") {
				continue
			}
			_k, _v, _ok := strings.Cut(_line, "=")
			if !_ok {
				continue
			}
			_k = strings.TrimSpace(_k)
			_v = strings.Trim(strings.TrimSpace(_v), `"'`)
			if _, _set := os.LookupEnv(_k); !_set {
				os.Setenv(_k, _v)
			}
		}
		_f.Close()
		return
	}
}
