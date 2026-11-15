package users

import "encoding/json"

type Prefs struct {
	Del      bool  `json:"del"`
	DelMe    bool  `json:"delMe"`
	DelMute  bool  `json:"delMute"`
	Edit     bool  `json:"edit"`
	EditMe   bool  `json:"editMe"`
	EditMute bool  `json:"editMute"`
	Media    bool  `json:"media"`
	Photo    bool  `json:"photo"`
	Video    bool  `json:"video"`
	Voice    bool  `json:"voice"`
	Files    bool  `json:"files"`
	Once     bool  `json:"once"`
	OnceAuto bool  `json:"onceAuto"`
	OnceHide bool  `json:"onceHide"`
	OnceMute bool  `json:"onceMute"`
	Panic    bool  `json:"panic"`
	PanicN   int   `json:"panicN"`
	Keep     int   `json:"keep"`
	Own      *bool `json:"own,omitempty"`
}

var __keeps__ = map[int]bool{0: true, 7: true, 30: true, 90: true, 365: true}

var __panics__ = map[int]bool{3: true, 5: true, 10: true, 20: true, 50: true}

func Def() Prefs {
	return Prefs{
		Del: true, Edit: true,
		Media: true, Photo: true, Video: true, Voice: true, Files: true,
		Once: true, OnceAuto: true, OnceHide: true,
		Panic: true, PanicN: 5,
		Keep: 30,
	}
}

func Parse(__raw []byte) Prefs {
	_p := Def()
	if len(__raw) > 0 {
		_ = json.Unmarshal(__raw, &_p)
	}
	return _p.Fix()
}

func (p Prefs) Fix() Prefs {
	if !__keeps__[p.Keep] {
		p.Keep = 30
	}
	if !__panics__[p.PanicN] {
		p.PanicN = 5
	}
	if p.Own != nil {
		p.DelMe, p.EditMe = *p.Own, *p.Own
		p.Own = nil
	}
	return p
}

func (p Prefs) Wants(__kind string) bool {
	if !p.Media {
		return false
	}
	switch __kind {
	case "photo":
		return p.Photo
	case "video", "round", "anim":
		return p.Video
	case "voice", "audio":
		return p.Voice
	}
	return p.Files
}

func (p Prefs) Mute(__kind string) bool {
	switch __kind {
	case "del":
		return p.DelMute
	case "edit":
		return p.EditMute
	case "once":
		return p.OnceMute
	}
	return false
}
