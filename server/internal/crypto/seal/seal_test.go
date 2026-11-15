package seal

import (
	"bytes"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"
)

func TestRound(t *testing.T) {
	_pub, _sec, _ := box.GenerateKey(rand.Reader)
	_msg := []byte("привет")
	_ct, _err := Seal(_pub[:], _msg)
	if _err != nil {
		t.Fatal(_err)
	}
	_out, _ok := box.OpenAnonymous(nil, _ct, _pub, _sec)
	if !_ok || !bytes.Equal(_out, _msg) {
		t.Fatal("open failed")
	}
}

func TestBadKey(t *testing.T) {
	if _, _err := Seal([]byte{1, 2}, []byte("x")); _err != ErrKey {
		t.Fatal("want ErrKey")
	}
}
