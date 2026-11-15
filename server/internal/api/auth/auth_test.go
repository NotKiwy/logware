package auth

import (
	"net/url"
	"strconv"
	"testing"
	"time"
)

const __tok__ = "123:abc"

func forge(__user string, __at time.Time) string {
	_v := url.Values{}
	_v.Set("auth_date", strconv.FormatInt(__at.Unix(), 10))
	_v.Set("query_id", "q1")
	_v.Set("user", __user)
	_v.Set("hash", Sign(__tok__, "auth_date="+_v.Get("auth_date")+"\nquery_id=q1\nuser="+__user))
	return _v.Encode()
}

func TestOk(t *testing.T) {
	_now := time.Now()
	_id, _err := Check(__tok__, forge(`{"id":42,"first_name":"A"}`, _now), _now)
	if _err != nil || _id != 42 {
		t.Fatalf("got %d %v", _id, _err)
	}
}

func TestForged(t *testing.T) {
	_now := time.Now()
	_raw := forge(`{"id":42}`, _now)
	if _, _err := Check("999:zzz", _raw, _now); _err != ErrBad {
		t.Fatal("want ErrBad for wrong token")
	}
}

func TestOld(t *testing.T) {
	_now := time.Now()
	if _, _err := Check(__tok__, forge(`{"id":42}`, _now.Add(-48*time.Hour)), _now); _err != ErrBad {
		t.Fatal("want ErrBad for old data")
	}
}
