package coin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"logware/internal/coin/cb"
	"logware/internal/coin/xr"
)

func TestCb(t *testing.T) {
	_tok := "123:abc"
	_body := []byte(`{"update_type":"invoice_paid"}`)
	_k := sha256.Sum256([]byte(_tok))
	_m := hmac.New(sha256.New, _k[:])
	_m.Write(_body)
	_sig := hex.EncodeToString(_m.Sum(nil))
	if !cb.Verify(_tok, _body, _sig) {
		t.Fatal("good")
	}
	if cb.Verify(_tok, append(_body, ' '), _sig) || cb.Verify("x", _body, _sig) || cb.Verify(_tok, _body, "zz") {
		t.Fatal("bad accepted")
	}
}

func TestXr(t *testing.T) {
	_sec := "whsec"
	_body := []byte(`{"type":"invoice"}`)
	_now := time.Now()
	_ts := strconv.FormatInt(_now.UnixMilli(), 10)
	_m := hmac.New(sha256.New, []byte(_sec))
	_m.Write([]byte(_ts + "." + string(_body)))
	_sig := hex.EncodeToString(_m.Sum(nil))
	if !xr.Verify(_sec, _body, _sig, "v1", _ts, _now) {
		t.Fatal("good")
	}
	if xr.Verify(_sec, _body, _sig, "v2", _ts, _now) || xr.Verify(_sec, _body, _sig, "v1", _ts, _now.Add(10*time.Minute)) || xr.Verify("x", _body, _sig, "v1", _ts, _now) {
		t.Fatal("bad accepted")
	}
}
