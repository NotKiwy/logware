package kw

import (
	"encoding/base64"
	"testing"
)

func TestHashes(t *testing.T) {
	_want := base64.StdEncoding.EncodeToString(Hash(42, "встреча завтра"))
	_found := false
	for _, _h := range Hashes(42, "Короче, ВСТРЕЧА — завтра в 7!") {
		if base64.StdEncoding.EncodeToString(_h) == _want {
			_found = true
		}
	}
	if !_found {
		t.Fatal("phrase not matched")
	}
	if got := base64.StdEncoding.EncodeToString(Hash(42, "ёлка")); got == base64.StdEncoding.EncodeToString(Hash(42, "елка")) {
		t.Fatal("hash must be on normalized words only")
	}
	if w := Words("Ёлка-палка"); len(w) != 2 || w[0] != "елка" {
		t.Fatal(w)
	}
}
