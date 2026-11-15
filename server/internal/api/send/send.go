package send

import (
	"html"
	"io"
	"net/http"
	"strconv"
	"unicode/utf8"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
	"logware/internal/tg/client"
)

type _Way struct {
	method string
	field  string
	cap    bool
}

var __ways__ = map[string]_Way{
	"photo":   {"sendPhoto", "photo", true},
	"video":   {"sendVideo", "video", true},
	"round":   {"sendVideoNote", "video_note", false},
	"voice":   {"sendVoice", "voice", true},
	"audio":   {"sendAudio", "audio", true},
	"anim":    {"sendAnimation", "animation", true},
	"sticker": {"sendSticker", "sticker", false},
	"doc":     {"sendDocument", "document", true},
}

const __capMax__ = 1024

func cut(__s string, __n int) string {
	if utf8.RuneCountInString(__s) <= __n {
		return __s
	}
	_r := []rune(__s)
	return string(_r[:__n-1]) + "…"
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, client.MaxFile+(1<<20))
		if _err := r.ParseMultipartForm(4 << 20); _err != nil {
			out.Fail(w, http.StatusBadRequest, "bad form")
			return
		}
		_uid := auth.ID(r)
		_chat := strconv.FormatInt(_uid, 10)
		_text := r.FormValue("text")
		if _fid := r.FormValue("id"); _fid != "" {
			_way, _ok := __ways__[r.FormValue("kind")]
			if !_ok {
				out.Fail(w, http.StatusBadRequest, "bad kind")
				return
			}
			_args := map[string]any{"chat_id": _uid, _way.field: _fid}
			if _way.cap && _text != "" {
				_args["caption"] = cut(_text, __capMax__)
			}
			if _err := __d.Tg.Call(r.Context(), _way.method, _args, nil); _err != nil {
				out.Oops(w, _err)
				return
			}
			out.JSON(w, map[string]bool{"ok": true})
			return
		}
		_f, _hdr, _err := r.FormFile("file")
		if _err != nil {
			if _text == "" {
				out.Fail(w, http.StatusBadRequest, "empty")
				return
			}
			if _err = __d.Tg.Send(r.Context(), _uid, html.EscapeString(cut(_text, 4096)), nil); _err != nil {
				out.Oops(w, _err)
				return
			}
			out.JSON(w, map[string]bool{"ok": true})
			return
		}
		defer _f.Close()
		_way, _ok := __ways__[r.FormValue("kind")]
		if !_ok {
			_way = __ways__["doc"]
		}
		_data, _err := io.ReadAll(_f)
		if _err != nil {
			out.Fail(w, http.StatusBadRequest, "bad file")
			return
		}
		_args := map[string]string{"chat_id": _chat}
		if _way.cap && _text != "" {
			_args["caption"] = cut(_text, __capMax__)
		}
		if _err = __d.Tg.Upload(r.Context(), _way.method, _way.field, _args, _hdr.Filename, _data); _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, map[string]bool{"ok": true})
	}
}
