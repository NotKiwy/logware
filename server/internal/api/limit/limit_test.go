package limit

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	_l := New(2, 3)
	_t := time.Unix(1000, 0)
	for _i := 0; _i < 3; _i++ {
		if !_l.Allow("a", _t) {
			t.Fatal("burst", _i)
		}
	}
	if _l.Allow("a", _t) {
		t.Fatal("over burst")
	}
	if !_l.Allow("b", _t) {
		t.Fatal("other key")
	}
	if !_l.Allow("a", _t.Add(600*time.Millisecond)) {
		t.Fatal("refill")
	}
	if _l.Allow("a", _t.Add(600*time.Millisecond)) {
		t.Fatal("refill once")
	}
	_l.Allow("x", _t)
	_l.Allow("y", _t.Add(2*time.Minute))
	if _, _ok := _l.all["x"]; _ok {
		t.Fatal("not swept")
	}
}
