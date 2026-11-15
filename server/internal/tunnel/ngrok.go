package tunnel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

var ErrNone = errors.New("PUBLIC_URL is empty and no ngrok tunnel found on 127.0.0.1:4040")

func Find(__ctx context.Context) (string, error) {
	_ctx, _done := context.WithTimeout(__ctx, 3*time.Second)
	defer _done()
	_req, _err := http.NewRequestWithContext(_ctx, http.MethodGet, "http://127.0.0.1:4040/api/tunnels", nil)
	if _err != nil {
		return "", _err
	}
	_res, _err := http.DefaultClient.Do(_req)
	if _err != nil {
		return "", ErrNone
	}
	defer _res.Body.Close()
	var _out struct {
		Tunnels []struct {
			URL string `json:"public_url"`
		} `json:"tunnels"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_out); _err != nil {
		return "", _err
	}
	for _, _t := range _out.Tunnels {
		if strings.HasPrefix(_t.URL, "https://") {
			return strings.TrimRight(_t.URL, "/"), nil
		}
	}
	return "", ErrNone
}
