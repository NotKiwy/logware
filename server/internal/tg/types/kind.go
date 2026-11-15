package types

import "encoding/json"

var __meta__ = map[string]bool{
	"message_id": true, "message_thread_id": true, "from": true, "sender_chat": true,
	"sender_boost_count": true, "sender_business_bot": true, "sender_tag": true, "date": true,
	"business_connection_id": true, "chat": true, "forward_origin": true, "is_topic_message": true,
	"is_automatic_forward": true, "reply_to_message": true, "external_reply": true, "quote": true,
	"reply_to_story": true, "reply_to_checklist_task_id": true, "via_bot": true, "edit_date": true,
	"has_protected_content": true, "is_from_offline": true, "is_paid_post": true, "media_group_id": true,
	"author_signature": true, "paid_star_count": true, "entities": true, "link_preview_options": true,
	"effect_id": true, "has_media_spoiler": true, "show_caption_above_media": true, "reply_markup": true,
	"caption_entities": true,
}

func (m *Message) Media() (*File, string) {
	switch {
	case len(m.Photo) > 0:
		return &m.Photo[len(m.Photo)-1], "photo"
	case m.Video != nil:
		return m.Video, "video"
	case m.Round != nil:
		return m.Round, "round"
	case m.Voice != nil:
		return m.Voice, "voice"
	case m.Anim != nil:
		return m.Anim, "anim"
	case m.Audio != nil:
		return m.Audio, "audio"
	case m.Doc != nil:
		return m.Doc, "doc"
	case m.Sticker != nil:
		return m.Sticker, "sticker"
	}
	return nil, ""
}

func (m *Message) Kind() string {
	if _, _k := m.Media(); _k != "" {
		return _k
	}
	if m.Text != "" {
		return "text"
	}
	if len(m.Rich) > 0 {
		return "rich"
	}
	var _raw map[string]json.RawMessage
	if json.Unmarshal(m.Raw, &_raw) != nil {
		return "empty"
	}
	for _k := range _raw {
		if !__meta__[_k] {
			return "other"
		}
	}
	return "empty"
}

func (m *Message) Out(__owner int64) bool {
	return m.From != nil && m.From.ID == __owner
}
