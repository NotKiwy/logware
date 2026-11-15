package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const __api__ = "https://api.telegram.org/"

const MaxFile = 20 << 20

type Err struct {
	Code int
	Desc string
}

func (e *Err) Error() string {
	return fmt.Sprintf("telegram %d: %s", e.Code, e.Desc)
}

type Client struct {
	__tok  string
	__http *http.Client
}

func New(__tok string) *Client {
	return &Client{__tok: __tok, __http: &http.Client{Timeout: 90 * time.Second}}
}

func (c *Client) Call(__ctx context.Context, __method string, __args any, __out any) error {
	_body, _err := json.Marshal(__args)
	if _err != nil {
		return _err
	}
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodPost, __api__+"bot"+c.__tok+"/"+__method, bytes.NewReader(_body))
	if _err != nil {
		return _err
	}
	_req.Header.Set("Content-Type", "application/json")
	_res, _err := c.__http.Do(_req)
	if _err != nil {
		return scrub(_err, c.__tok)
	}
	defer _res.Body.Close()
	var _r struct {
		Ok   bool            `json:"ok"`
		Res  json.RawMessage `json:"result"`
		Desc string          `json:"description"`
		Code int             `json:"error_code"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_r); _err != nil {
		return _err
	}
	if !_r.Ok {
		return &Err{Code: _r.Code, Desc: _r.Desc}
	}
	if __out == nil {
		return nil
	}
	return json.Unmarshal(_r.Res, __out)
}

func (c *Client) Fetch(__ctx context.Context, __path string) ([]byte, error) {
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodGet, __api__+"file/bot"+c.__tok+"/"+__path, nil)
	if _err != nil {
		return nil, _err
	}
	_res, _err := c.__http.Do(_req)
	if _err != nil {
		return nil, scrub(_err, c.__tok)
	}
	defer _res.Body.Close()
	if _res.StatusCode != http.StatusOK {
		return nil, &Err{Code: _res.StatusCode, Desc: "file download failed"}
	}
	return io.ReadAll(io.LimitReader(_res.Body, MaxFile+1))
}

func scrub(__err error, __tok string) error {
	return errors.New(string(bytes.ReplaceAll([]byte(__err.Error()), []byte(__tok), []byte("***"))))
}
