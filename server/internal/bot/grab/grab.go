package grab

import (
	"context"
	"errors"

	"logware/internal/crypto/seal"
	"logware/internal/deps"
	"logware/internal/tg/client"
	"logware/internal/tg/types"
)

var ErrBig = errors.New("file is larger than 20 MB")

func Media(__ctx context.Context, __d *deps.D, __pub []byte, __f *types.File) ([]byte, error) {
	if __f.Size > client.MaxFile {
		return nil, ErrBig
	}
	_f, _err := __d.Tg.File(__ctx, __f.ID)
	if _err != nil {
		return nil, _err
	}
	if _f.Path == "" {
		return nil, ErrBig
	}
	_b, _err := __d.Tg.Fetch(__ctx, _f.Path)
	if _err != nil {
		return nil, _err
	}
	if len(_b) > client.MaxFile {
		return nil, ErrBig
	}
	return seal.Seal(__pub, _b)
}
