package kw

import (
	"context"
	"crypto/sha256"
	"strconv"
	"strings"
	"unicode"

	"logware/internal/bot/notify"
	"logware/internal/bot/own"
	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/types"
)

const __span__ = 3

func Words(__s string) []string {
	_s := strings.ReplaceAll(strings.ToLower(__s), "ё", "е")
	return strings.FieldsFunc(_s, func(_r rune) bool { return !unicode.IsLetter(_r) && !unicode.IsNumber(_r) })
}

func Hash(__user int64, __phrase string) []byte {
	_h := sha256.Sum256([]byte("lw:" + strconv.FormatInt(__user, 10) + ":" + __phrase))
	return _h[:]
}

func Hashes(__user int64, __s string) [][]byte {
	_ws := Words(__s)
	_seen := map[string]bool{}
	var _out [][]byte
	for _i := range _ws {
		for _n := 1; _n <= __span__ && _i+_n <= len(_ws); _n++ {
			_p := strings.Join(_ws[_i:_i+_n], " ")
			if _seen[_p] {
				continue
			}
			_seen[_p] = true
			_out = append(_out, Hash(__user, _p))
		}
	}
	return _out
}

func Check(__ctx context.Context, __d *deps.D, __o *own.Own, __m *types.Message, __row int64) error {
	_t := __m.Text
	if _t == "" {
		_t = __m.Caption
	}
	if _t == "" || __row == 0 {
		return nil
	}
	_hit, _err := __d.Kws.Hit(__ctx, __o.User.ID, Hashes(__o.User.ID, _t))
	if _err != nil || !_hit {
		return _err
	}
	return notify.Send(__ctx, __d, __o, "kw", []int64{__row}, text.Kw(text.Name(__m.Chat)))
}
