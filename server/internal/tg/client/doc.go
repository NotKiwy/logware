package client

import (
	"context"

	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

func (c *Client) Post(__ctx context.Context, __chat int64, __text string, __kb *kb.Kb) (int64, error) {
	_args := map[string]any{
		"chat_id":              __chat,
		"text":                 __text,
		"parse_mode":           "HTML",
		"link_preview_options": map[string]bool{"is_disabled": true},
	}
	if __kb != nil {
		_args["reply_markup"] = __kb
	}
	var _m types.Message
	return _m.ID, c.Call(__ctx, "sendMessage", _args, &_m)
}

func (c *Client) Answer(__ctx context.Context, __id, __text string) error {
	return c.Call(__ctx, "answerCallbackQuery", map[string]any{"callback_query_id": __id, "text": __text}, nil)
}

func (c *Client) Retext(__ctx context.Context, __chat, __msg int64, __text string) error {
	return c.Call(__ctx, "editMessageText", map[string]any{
		"chat_id": __chat, "message_id": __msg, "text": __text, "parse_mode": "HTML",
		"link_preview_options": map[string]bool{"is_disabled": true},
		"reply_markup":         map[string]any{"inline_keyboard": [][]kb.Btn{}},
	}, nil)
}
