package cb

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

const __api__ = "https://pay.crypt.bot/api/"

var __http = &http.Client{Timeout: 20 * time.Second}

type Inv struct {
	ID      int64  `json:"invoice_id"`
	Status  string `json:"status"`
	Bot     string `json:"bot_invoice_url"`
	App     string `json:"mini_app_invoice_url"`
	Payload string `json:"payload"`
	Fiat    string `json:"fiat"`
	Amount  string `json:"amount"`
	Asset   string `json:"paid_asset"`
	Paid    string `json:"paid_amount"`
}

func (i *Inv) Link() string {
	if i.App != "" {
		return i.App
	}
	return i.Bot
}

type Update struct {
	ID   int64  `json:"update_id"`
	Type string `json:"update_type"`
	Inv  Inv    `json:"payload"`
}

func Create(__ctx context.Context, __tok string, __rub int, __desc, __payload, __back string) (*Inv, error) {
	_body, _ := json.Marshal(map[string]any{
		"currency_type": "fiat",
		"fiat":          "RUB",
		"amount":        strconv.Itoa(__rub),
		"description":   __desc,
		"payload":       __payload,
		"expires_in":    3600,
		"paid_btn_name": "openBot",
		"paid_btn_url":  __back,
	})
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodPost, __api__+"createInvoice", bytes.NewReader(_body))
	if _err != nil {
		return nil, _err
	}
	_req.Header.Set("Content-Type", "application/json")
	_req.Header.Set("Crypto-Pay-API-Token", __tok)
	_res, _err := __http.Do(_req)
	if _err != nil {
		return nil, errors.New("cryptobot: request failed")
	}
	defer _res.Body.Close()
	var _r struct {
		Ok     bool `json:"ok"`
		Result Inv  `json:"result"`
		Error  struct {
			Name string `json:"name"`
		} `json:"error"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_r); _err != nil {
		return nil, _err
	}
	if !_r.Ok {
		return nil, errors.New("cryptobot: " + _r.Error.Name)
	}
	return &_r.Result, nil
}

func Verify(__tok string, __body []byte, __sig string) bool {
	_got, _err := hex.DecodeString(__sig)
	if _err != nil || __tok == "" {
		return false
	}
	_key := sha256.Sum256([]byte(__tok))
	_m := hmac.New(sha256.New, _key[:])
	_m.Write(__body)
	return hmac.Equal(_m.Sum(nil), _got)
}

func Paid(__ctx context.Context, __tok string) ([]Inv, error) {
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodGet, __api__+"getInvoices?status=paid&count=100", nil)
	if _err != nil {
		return nil, _err
	}
	_req.Header.Set("Crypto-Pay-API-Token", __tok)
	_res, _err := __http.Do(_req)
	if _err != nil {
		return nil, errors.New("cryptobot: request failed")
	}
	defer _res.Body.Close()
	var _r struct {
		Ok     bool `json:"ok"`
		Result struct {
			Items []Inv `json:"items"`
		} `json:"result"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_r); _err != nil {
		return nil, _err
	}
	if !_r.Ok {
		return nil, errors.New("cryptobot: getInvoices failed")
	}
	return _r.Result.Items, nil
}
