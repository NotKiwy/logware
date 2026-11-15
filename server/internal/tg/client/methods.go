package client

import (
	"context"

	"logware/internal/tg/kb"
	"logware/internal/tg/types"
)

var __allowed__ = []string{
	"message", "callback_query", "business_connection",
	"business_message", "edited_business_message", "deleted_business_messages", "pre_checkout_query",
}

func (c *Client) Me(__ctx context.Context) (*types.User, error) {
	var _u types.User
	return &_u, c.Call(__ctx, "getMe", struct{}{}, &_u)
}

func (c *Client) NoHook(__ctx context.Context) error {
	return c.Call(__ctx, "deleteWebhook", struct{}{}, nil)
}

func (c *Client) Updates(__ctx context.Context, __off int64) ([]types.Update, error) {
	var _u []types.Update
	_err := c.Call(__ctx, "getUpdates", map[string]any{
		"offset":          __off,
		"timeout":         50,
		"allowed_updates": __allowed__,
	}, &_u)
	return _u, _err
}

func (c *Client) Send(__ctx context.Context, __chat int64, __text string, __kb *kb.Kb) error {
	return c.Note(__ctx, __chat, __text, __kb, false)
}

func (c *Client) Note(__ctx context.Context, __chat int64, __text string, __kb *kb.Kb, __mute bool) error {
	_args := map[string]any{
		"chat_id":              __chat,
		"text":                 __text,
		"parse_mode":           "HTML",
		"link_preview_options": map[string]bool{"is_disabled": true},
		"disable_notification": __mute,
	}
	if __kb != nil {
		_args["reply_markup"] = __kb
	}
	return c.Call(__ctx, "sendMessage", _args, nil)
}

func (c *Client) File(__ctx context.Context, __id string) (*types.File, error) {
	var _f types.File
	return &_f, c.Call(__ctx, "getFile", map[string]string{"file_id": __id}, &_f)
}

func (c *Client) Conn(__ctx context.Context, __id string) (*types.Conn, error) {
	var _c types.Conn
	return &_c, c.Call(__ctx, "getBusinessConnection", map[string]string{"business_connection_id": __id}, &_c)
}

func (c *Client) DelBiz(__ctx context.Context, __conn string, __ids []int64) error {
	return c.Call(__ctx, "deleteBusinessMessages", map[string]any{
		"business_connection_id": __conn,
		"message_ids":            __ids,
	}, nil)
}

func (c *Client) Menu(__ctx context.Context, __text, __url string) error {
	return c.Call(__ctx, "setChatMenuButton", map[string]any{
		"menu_button": map[string]any{"type": "web_app", "text": __text, "web_app": kb.App{URL: __url}},
	}, nil)
}

func (c *Client) SendBiz(__ctx context.Context, __conn string, __chat int64, __text string, __reply int64) error {
	return c.Call(__ctx, "sendMessage", map[string]any{
		"business_connection_id": __conn,
		"chat_id":                __chat,
		"text":                   __text,
		"parse_mode":             "HTML",
		"reply_parameters":       map[string]any{"message_id": __reply, "allow_sending_without_reply": true},
	}, nil)
}

func (c *Client) BizName(__ctx context.Context, __conn, __first, __last string) error {
	return c.Call(__ctx, "setBusinessAccountName", map[string]string{
		"business_connection_id": __conn,
		"first_name":             __first,
		"last_name":              __last,
	}, nil)
}

func (c *Client) Photo(__ctx context.Context, __user int64) (*types.File, error) {
	var _r struct {
		Photos [][]types.File `json:"photos"`
	}
	if _err := c.Call(__ctx, "getUserProfilePhotos", map[string]any{"user_id": __user, "limit": 1}, &_r); _err != nil {
		return nil, _err
	}
	if len(_r.Photos) == 0 || len(_r.Photos[0]) == 0 {
		return nil, nil
	}
	_s := _r.Photos[0]
	_pick := _s[0]
	if len(_s) > 1 {
		_pick = _s[1]
	}
	return c.File(__ctx, _pick.ID)
}

func (c *Client) Emojis(__ctx context.Context, __ids []string) ([]types.Emoji, error) {
	var _out []types.Emoji
	return _out, c.Call(__ctx, "getCustomEmojiStickers", map[string]any{"custom_emoji_ids": __ids}, &_out)
}

func (c *Client) EditBiz(__ctx context.Context, __conn string, __chat, __msg int64, __text string) error {
	return c.Call(__ctx, "editMessageText", map[string]any{
		"business_connection_id": __conn,
		"chat_id":                __chat,
		"message_id":             __msg,
		"text":                   __text,
		"parse_mode":             "HTML",
	}, nil)
}

func (c *Client) Cmds(__ctx context.Context, __cmds [][2]string) error {
	_l := make([]map[string]string, 0, len(__cmds))
	for _, _c := range __cmds {
		_l = append(_l, map[string]string{"command": _c[0], "description": _c[1]})
	}
	return c.Call(__ctx, "setMyCommands", map[string]any{"commands": _l}, nil)
}
