package types

import "encoding/json"

type Update struct {
	ID      int64     `json:"update_id"`
	Msg     *Message  `json:"message"`
	Cb      *Callback `json:"callback_query"`
	Conn    *Conn     `json:"business_connection"`
	BizMsg  *Message  `json:"business_message"`
	BizEdit *Message  `json:"edited_business_message"`
	BizDel  *Deleted  `json:"deleted_business_messages"`
	Pre     *Pre      `json:"pre_checkout_query"`
}

type Pre struct {
	ID       string `json:"id"`
	From     User   `json:"from"`
	Currency string `json:"currency"`
	Total    int    `json:"total_amount"`
	Payload  string `json:"invoice_payload"`
}

type Paid struct {
	Currency  string `json:"currency"`
	Total     int    `json:"total_amount"`
	Payload   string `json:"invoice_payload"`
	Charge    string `json:"telegram_payment_charge_id"`
	Expires   int64  `json:"subscription_expiration_date"`
	Recurring bool   `json:"is_recurring"`
}

type User struct {
	ID    int64  `json:"id"`
	Bot   bool   `json:"is_bot"`
	First string `json:"first_name"`
	Last  string `json:"last_name"`
	Nick  string `json:"username"`
	Prem  bool   `json:"is_premium"`
	Lang  string `json:"language_code"`
	Biz   bool   `json:"can_connect_to_business"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	First string `json:"first_name"`
	Last  string `json:"last_name"`
	Nick  string `json:"username"`
	Title string `json:"title"`
}

type File struct {
	ID   string `json:"file_id"`
	Uniq string `json:"file_unique_id"`
	Size int64  `json:"file_size"`
	Path string `json:"file_path"`
	Mime string `json:"mime_type"`
	Name string `json:"file_name"`
	Dur  int    `json:"duration"`
}

type Rights struct {
	Reply  bool `json:"can_reply"`
	Read   bool `json:"can_read_messages"`
	DelOwn bool `json:"can_delete_sent_messages"`
	DelAll bool `json:"can_delete_all_messages"`
	Name   bool `json:"can_edit_name"`
}

type Conn struct {
	ID     string `json:"id"`
	User   User   `json:"user"`
	ChatID int64  `json:"user_chat_id"`
	Date   int64  `json:"date"`
	Rights Rights `json:"rights"`
	On     bool   `json:"is_enabled"`
}

type Deleted struct {
	Conn string  `json:"business_connection_id"`
	Chat Chat    `json:"chat"`
	IDs  []int64 `json:"message_ids"`
}

type Callback struct {
	ID   string   `json:"id"`
	From User     `json:"from"`
	Msg  *Message `json:"message"`
	Data string   `json:"data"`
}

type Message struct {
	ID      int64           `json:"message_id"`
	From    *User           `json:"from"`
	Chat    Chat            `json:"chat"`
	Date    int64           `json:"date"`
	Conn    string          `json:"business_connection_id"`
	ByBot   *User           `json:"sender_business_bot"`
	Reply   *Message        `json:"reply_to_message"`
	Text    string          `json:"text"`
	Caption string          `json:"caption"`
	Photo   []File          `json:"photo"`
	Video   *File           `json:"video"`
	Anim    *File           `json:"animation"`
	Audio   *File           `json:"audio"`
	Voice   *File           `json:"voice"`
	Round   *File           `json:"video_note"`
	Doc     *File           `json:"document"`
	Sticker *File           `json:"sticker"`
	Rich    json.RawMessage `json:"rich_message"`
	Ents    []Ent           `json:"entities"`
	Paid    *Paid           `json:"successful_payment"`
	Raw     json.RawMessage `json:"-"`
}

type Ent struct {
	Type  string `json:"type"`
	Emoji string `json:"custom_emoji_id"`
}

type Emoji struct {
	ID    string `json:"custom_emoji_id"`
	Emoji string `json:"emoji"`
	Set   string `json:"set_name"`
}

func (m *Message) UnmarshalJSON(__b []byte) error {
	type _m Message
	if _err := json.Unmarshal(__b, (*_m)(m)); _err != nil {
		return _err
	}
	m.Raw = append(json.RawMessage(nil), __b...)
	return nil
}
