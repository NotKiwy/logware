package seal

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/nacl/box"
)

const KeyLen = 32

var ErrKey = errors.New("bad public key")

func Seal(__pub, __msg []byte) ([]byte, error) {
	if len(__pub) != KeyLen {
		return nil, ErrKey
	}
	var _k [KeyLen]byte
	copy(_k[:], __pub)
	return box.SealAnonymous(nil, __msg, &_k, rand.Reader)
}
