package clock

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"logware/internal/mods/mod"
)

const __max__ = 64

var __days__ = [...]string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

var __faces__ = [...]string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}

var __halves__ = [...]string{"🕧", "🕜", "🕝", "🕞", "🕟", "🕠", "🕡", "🕢", "🕣", "🕤", "🕥", "🕦"}

func now(__c mod.Cfg) time.Time {
	_loc, _err := time.LoadLocation(__c["tz"])
	if _err != nil {
		_loc, _ = time.LoadLocation("Europe/Moscow")
	}
	return time.Now().In(_loc)
}

var __full__ = [...]string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}

var __mons__ = [...]string{"янв", "фев", "мар", "апр", "мая", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}

var __months__ = [...]string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}

func Stamp(__tpl string, __t time.Time) string {
	_face := __faces__[__t.Hour()%12]
	if __t.Minute() >= 30 {
		_face = __halves__[__t.Hour()%12]
	}
	_out := strings.NewReplacer(
		"%DS", __days__[__t.Weekday()],
		"%BS", __mons__[__t.Month()-1],
		"%D", __full__[__t.Weekday()],
		"%B", __months__[__t.Month()-1],
		"%H", __t.Format("15"),
		"%h", __t.Format("3"),
		"%M", __t.Format("04"),
		"%P", __t.Format("PM"),
		"%d", __t.Format("02"),
		"%m", __t.Format("01"),
		"%y", __t.Format("06"),
		"%Y", __t.Format("2006"),
		"%E", _face,
	).Replace(__tpl)
	if strings.TrimSpace(_out) == "" {
		return __t.Format("15:04")
	}
	return cut(_out, 40)
}

func cut(__s string, __n int) string {
	if utf8.RuneCountInString(__s) <= __n {
		return __s
	}
	return string([]rune(__s)[:__n])
}

func Name(__c mod.Cfg, __t time.Time) (string, string) {
	_tag := Stamp(__c["tpl"], __t)
	_first, _last := strings.TrimSpace(__c["first"]), strings.TrimSpace(__c["last"])
	_room := __max__ - utf8.RuneCountInString(_tag) - 1
	switch {
	case __c["pos"] == "before":
		_first = _tag + " " + cut(_first, _room)
	case _last != "":
		_last = cut(_last, _room) + " " + _tag
	default:
		_first = cut(_first, _room) + " " + _tag
	}
	return strings.TrimSpace(_first), strings.TrimSpace(_last)
}

func Key(__f, __l string) string {
	return fmt.Sprintf("%s\x00%s", __f, __l)
}
