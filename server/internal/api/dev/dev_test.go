package dev

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"
)

func TestRepo(t *testing.T) {
	_ok := map[string]string{
		"https://github.com/NotKiwy/logware":  "https://github.com/NotKiwy/logware",
		"github.com/a-b/c.d_e.git":            "https://github.com/a-b/c.d_e",
		" https://www.github.com/user/repo/ ": "https://github.com/user/repo",
	}
	for _in, _want := range _ok {
		if _got, _good := Repo(_in); !_good || _got != _want {
			t.Fatal(_in, _got)
		}
	}
	for _, _bad := range []string{"", "https://gitlab.com/a/b", "github.com/a", "github.com/a/b/c", "javascript:alert(1)", "https://github.com/a/b?x=1"} {
		if _, _good := Repo(_bad); _good {
			t.Fatal(_bad)
		}
	}
}

func png64(__w, __h int) string {
	var _b bytes.Buffer
	_ = png.Encode(&_b, image.NewRGBA(image.Rect(0, 0, __w, __h)))
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(_b.Bytes())
}

func TestPic(t *testing.T) {
	if _b, _ok := Pic(png64(256, 256)); !_ok || len(_b) == 0 {
		t.Fatal("good png")
	}
	if _b, _ok := Pic(""); !_ok || _b != nil {
		t.Fatal("empty")
	}
	for _, _bad := range []string{png64(10, 10), png64(2000, 2000), "data:image/png;base64,AAAA", "not a pic"} {
		if _, _ok := Pic(_bad); _ok {
			t.Fatal("bad accepted")
		}
	}
}
