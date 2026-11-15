package help

import (
	"context"
	"strings"

	"logware/internal/bot/text"
	"logware/internal/deps"
	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

const __docs__ = "https://docs.logware.dev"

func Run(__ctx context.Context, __d *deps.D, __m *types.Message) (bool, error) {
	_cmd, _, _ := strings.Cut(strings.TrimSpace(__m.Text), " ")
	_cmd, _, _ = strings.Cut(_cmd, "@")
	switch _cmd {
	case "/terms":
		return true, __d.Tg.Send(__ctx, __m.Chat.ID, text.Terms(), kb.One(kb.Url(text.TermsBtn, __docs__+"/legal/terms")))
	case "/privacy":
		return true, __d.Tg.Send(__ctx, __m.Chat.ID, text.Privacy(), kb.One(kb.Url(text.PrivacyBtn, __docs__+"/legal/privacy")))
	case "/paysupport", "/support":
		if __d.Sup == "" {
			return true, __d.Tg.Send(__ctx, __m.Chat.ID, text.NoSupport(), nil)
		}
		return true, __d.Tg.Send(__ctx, __m.Chat.ID, text.Support(__d.Sup), kb.One(kb.Url(text.SupportBtn, "https://t.me/"+__d.Sup).With(text.Bell)))
	}
	return false, nil
}
