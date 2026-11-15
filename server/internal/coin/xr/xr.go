package xr

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
	"strings"
	"time"
)

const __api__ = "https://pay.api.xrocket.exchange/api/v1/"

const Window = 5 * time.Minute

var __http = &http.Client{Timeout: 20 * time.Second}

type Inv struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ClientID string `json:"clientInvoiceId"`
	Amount   string `json:"priceAmount"`
	Currency string `json:"priceCurrency"`
	Links    struct {
		Bot string `json:"telegramBotLink"`
	} `json:"links"`
	Customer struct {
		Tg string `json:"telegramId"`
	} `json:"customer"`
}

type Event struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Event   string `json:"event"`
		Inv     Inv    `json:"invoice"`
		Payment struct {
			Amount   string `json:"receiveAmount"`
			Currency string `json:"receiveCurrency"`
		} `json:"payment"`
	} `json:"data"`
}

func Create(__ctx context.Context, __tok string, __usdt, __desc, __client string, __user int64, __back string) (*Inv, error) {
	_body, _ := json.Marshal(map[string]any{
		"priceAmount":     __usdt,
		"priceCurrency":   "USDT",
		"numPayments":     1,
		"clientInvoiceId": __client,
		"description":     __desc,
		"expiresIn":       3600000,
		"customer":        map[string]string{"telegramId": strconv.FormatInt(__user, 10)},
		"url":             map[string]string{"successUrl": __back, "cancelUrl": __back},
	})
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodPost, __api__+"invoices", bytes.NewReader(_body))
	if _err != nil {
		return nil, _err
	}
	_req.Header.Set("Content-Type", "application/json")
	_req.Header.Set("Authorization", "Bearer "+__tok)
	_res, _err := __http.Do(_req)
	if _err != nil {
		return nil, errors.New("xrocket: request failed")
	}
	defer _res.Body.Close()
	if _res.StatusCode >= 300 {
		var _p struct {
			Title  string `json:"title"`
			Detail string `json:"detail"`
		}
		_ = json.NewDecoder(_res.Body).Decode(&_p)
		return nil, errors.New("xrocket: " + strconv.Itoa(_res.StatusCode) + " " + _p.Title + " " + _p.Detail)
	}
	var _i Inv
	if _err = json.NewDecoder(_res.Body).Decode(&_i); _err != nil {
		return nil, _err
	}
	if _i.Links.Bot == "" {
		return nil, errors.New("xrocket: no link")
	}
	return &_i, nil
}

func Verify(__secret string, __body []byte, __sig, __ver, __ts string, __now time.Time) bool {
	if __secret == "" || __ver != "v1" || __sig == "" {
		return false
	}
	_ms, _err := strconv.ParseInt(__ts, 10, 64)
	if _err != nil {
		return false
	}
	_d := __now.Sub(time.UnixMilli(_ms))
	if _d > Window || _d < -Window {
		return false
	}
	_got, _err := hex.DecodeString(__sig)
	if _err != nil {
		return false
	}
	_m := hmac.New(sha256.New, []byte(__secret))
	_m.Write([]byte(__ts + "."))
	_m.Write(__body)
	return hmac.Equal(_m.Sum(nil), _got)
}

func Paid(__ctx context.Context, __tok string) ([]Inv, error) {
	_req, _err := http.NewRequestWithContext(__ctx, http.MethodGet, __api__+"invoices?status=paid&limit=100", nil)
	if _err != nil {
		return nil, _err
	}
	_req.Header.Set("Authorization", "Bearer "+__tok)
	_res, _err := __http.Do(_req)
	if _err != nil {
		return nil, errors.New("xrocket: request failed")
	}
	defer _res.Body.Close()
	if _res.StatusCode >= 300 {
		return nil, errors.New("xrocket: list " + strconv.Itoa(_res.StatusCode))
	}
	var _r struct {
		Items []Inv `json:"items"`
	}
	if _err = json.NewDecoder(_res.Body).Decode(&_r); _err != nil {
		return nil, _err
	}
	return _r.Items, nil
}

func Client(__id string) string {
	if strings.Count(__id, ":") == 2 {
		return __id[:strings.LastIndex(__id, ":")]
	}
	return __id
}
