package clock

import (
	"testing"
	"time"

	"logware/internal/mods/mod"
)

func TestName(t *testing.T) {
	_t := time.Date(2026, 11, 12, 23, 18, 0, 0, time.UTC)
	for _, _c := range []struct {
		cfg   mod.Cfg
		first string
		last  string
	}{
		{mod.Cfg{"tpl": "[%H:%M]", "pos": "after", "first": "Алекс", "last": "Журавлёв"}, "Алекс", "Журавлёв [23:18]"},
		{mod.Cfg{"tpl": "%h:%M %P |", "pos": "before", "first": "Алекс"}, "11:18 PM | Алекс", ""},
		{mod.Cfg{"tpl": "(%d.%m %H:%M)", "pos": "after", "first": "Алекс"}, "Алекс (12.11 23:18)", ""},
		{mod.Cfg{"tpl": "%DS, %d %BS %Y", "first": "A"}, "A Чт, 12 ноя 2026", ""},
		{mod.Cfg{"tpl": "%D %d %B %y", "first": "A"}, "A четверг 12 ноября 26", ""},
		{mod.Cfg{"tpl": "%E", "pos": "before", "first": "A"}, "🕚 A", ""},
		{mod.Cfg{"tpl": "  ", "first": "A"}, "A 23:18", ""},
	} {
		_f, _l := Name(_c.cfg, _t)
		if _f != _c.first || _l != _c.last {
			t.Errorf("%v: got %q %q", _c.cfg, _f, _l)
		}
	}
}
