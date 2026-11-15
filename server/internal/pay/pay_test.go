package pay

import (
	"testing"
	"time"
)

func TestPayload(t *testing.T) {
	_p, _u, _ok := Parse(Payload("pro", 42))
	if !_ok || _p.ID != "pro" || _u != 42 || _p.Price != 149 {
		t.Fatal(_p, _u, _ok)
	}
	if _p, _u, _ok := Parse(Tagged("dev", 7, "app")); !_ok || _p.ID != "dev" || _u != 7 {
		t.Fatal("tagged")
	}
	if _, _t := Tag(Tagged("dev", 7, "app")); _t != "app" {
		t.Fatal("tag")
	}
	for _, _s := range []string{"", "pro", "x:1", "pro:abc"} {
		if _, _, _ok := Parse(_s); _ok {
			t.Fatal(_s)
		}
	}
}

func TestMerge(t *testing.T) {
	_now := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)
	_d := 24 * time.Hour
	for _, _c := range []struct {
		plan  string
		left  time.Duration
		nw    string
		add   time.Duration
		want  string
		after time.Duration
	}{
		{"", 0, "pro", 30 * _d, "pro", 30 * _d},
		{"pro", -_d, "dev", 3 * _d, "dev", 3 * _d},
		{"pro", 10 * _d, "pro", 30 * _d, "pro", 40 * _d},
		{"pro", 365 * _d, "dev", 3 * _d, "dev", 3*_d + time.Duration(float64(365*_d)*149/299)},
		{"dev", 10 * _d, "pro", 30 * _d, "dev", 10*_d + time.Duration(float64(30*_d)*149/299)},
	} {
		_p, _u := Merge(_c.plan, _now.Add(_c.left), _c.nw, _c.add, _now)
		if _p != _c.want || _u.Sub(_now) != _c.after {
			t.Errorf("%+v: %s %v", _c, _p, _u.Sub(_now))
		}
	}
}
