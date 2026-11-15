package text

import (
	"strings"
	"testing"

	"logware/internal/tg/types"
)

func TestTexts(t *testing.T) {
	for _, _s := range []string{
		Hello(), Ready(true, "bot"), Ready(false, "bot"), Linked(types.Rights{}, true, "bot"), Unlinked(),
		Del("Маша", 1), Del("Маша", 3), Del("Маша", 120), Edit("Маша"), Kw("Маша"), Panic("Маша", 7), Saved("Маша"),
	} {
		if strings.Contains(_s, ".s<") || strings.Contains(_s, "<code>.s") {
			t.Errorf("hint left: %s", _s)
		}
		if strings.Count(_s, "<blockquote") != strings.Count(_s, "</blockquote>") || strings.Count(_s, "<b>") != strings.Count(_s, "</b>") {
			t.Errorf("unbalanced: %s", _s)
		}
		t.Log(_s)
	}
}
