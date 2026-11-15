package client

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

func (c *Client) Upload(__ctx context.Context, __method, __field string, __args map[string]string, __name string, __data []byte) error {
	var _buf bytes.Buffer
	_w := multipart.NewWriter(&_buf)
	for _k, _v := range __args {
		if _err := _w.WriteField(_k, _v); _err != nil {
			return _err
		}
	}
	_f, _err := _w.CreateFormFile(__field, __name)
	if _err != nil {
		return _err
	}
	if _, _err = _f.Write(__data); _err != nil {
		return _err
	}
	if _err = _w.Close(); _err != nil {
		return _err
	}
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodPost, __api__+"bot"+c.__tok+"/"+__method, &_buf)
	if _err != nil {
		return _err
	}
	_req.Header.Set("Content-Type", _w.FormDataContentType())
	_res, _err := c.__http.Do(_req)
	if _err != nil {
		return scrub(_err, c.__tok)
	}
	defer _res.Body.Close()
	var _r struct {
		Ok   bool   `json:"ok"`
		Desc string `json:"description"`
		Code int    `json:"error_code"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_r); _err != nil {
		return _err
	}
	if !_r.Ok {
		return &Err{Code: _r.Code, Desc: _r.Desc}
	}
	return nil
}
