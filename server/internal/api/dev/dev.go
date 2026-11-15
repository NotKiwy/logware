package dev

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"logware/internal/api/auth"
	"logware/internal/api/mods"
	"logware/internal/api/out"
	"logware/internal/bot/review"
	"logware/internal/deps"
	"logware/internal/store/devs"
)

var __nick__ = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,23}$`)

var __slug__ = regexp.MustCompile(`^[a-z][a-z0-9-]{1,23}$`)

var __repo__ = regexp.MustCompile(`^github\.com/[A-Za-z0-9](?:[A-Za-z0-9-]{0,38})/[A-Za-z0-9._-]{1,100}$`)

const __pic__ = 400 << 10

type Plug struct {
	ID    int64  `json:"id"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	About string `json:"about"`
	Ver   string `json:"ver"`
	Repo  string `json:"repo"`
	Live  bool   `json:"live"`
	Last  string `json:"last"`
	Why   string `json:"why"`
	Kind  string `json:"kind"`
	Users int    `json:"users"`
	Pic   string `json:"pic"`
}

type Req struct {
	ID      int64  `json:"id"`
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	About   string `json:"about"`
	Ver     string `json:"ver"`
	Note    string `json:"note"`
	Repo    string `json:"repo"`
	Pic     bool   `json:"pic"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Created int64  `json:"created"`
	Done    int64  `json:"done"`
}

type View struct {
	Nick  string `json:"nick"`
	Bio   string `json:"bio"`
	Plugs []Plug `json:"plugs"`
	Gh    string `json:"gh"`
}

type One struct {
	Plug
	Reqs []Req `json:"reqs"`
}

type In struct {
	Nick string `json:"nick"`
	Bio  string `json:"bio"`
}

type PlugIn struct {
	Kind  string `json:"kind"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	About string `json:"about"`
	Ver   string `json:"ver"`
	Note  string `json:"note"`
	Repo  string `json:"repo"`
	Pic   string `json:"pic"`
}

func can(__d *deps.D, r *http.Request) bool {
	_id := auth.ID(r)
	if _id == __d.Cfg.Own && _id != 0 {
		return true
	}
	_s, _err := __d.Subs.Get(r.Context(), _id)
	return _err == nil && _s.Live() && _s.Plan == "dev"
}

func Repo(__s string) (string, bool) {
	_s := strings.TrimSpace(__s)
	_s = strings.TrimPrefix(strings.TrimPrefix(_s, "https://"), "http://")
	_s = strings.TrimPrefix(_s, "www.")
	_s = strings.TrimSuffix(strings.TrimSuffix(_s, "/"), ".git")
	if !__repo__.MatchString(_s) {
		return "", false
	}
	return "https://" + _s, true
}

func Pic(__s string) ([]byte, bool) {
	if __s == "" {
		return nil, true
	}
	_, _b64, _ok := strings.Cut(__s, ";base64,")
	if !_ok {
		return nil, false
	}
	_b, _err := base64.StdEncoding.DecodeString(_b64)
	if _err != nil || len(_b) == 0 || len(_b) > __pic__ {
		return nil, false
	}
	_c, _f, _err := image.DecodeConfig(bytes.NewReader(_b))
	if _err != nil || (_f != "png" && _f != "jpeg") || _c.Width < 64 || _c.Height < 64 || _c.Width > 1024 || _c.Height > 1024 {
		return nil, false
	}
	return _b, true
}

func short(__s string, __max int) (string, bool) {
	_v := strings.TrimSpace(__s)
	return _v, utf8.RuneCountInString(_v) <= __max
}

func view(__p devs.Plug, __users int) Plug {
	return Plug{ID: __p.ID, Slug: __p.Slug, Name: __p.Name, About: __p.About, Ver: __p.Ver, Repo: __p.Repo, Live: __p.Live,
		Last: __p.Last, Why: __p.Why, Kind: __p.Kind, Users: __users, Pic: mods.PicURL(__p.ID, __p.Pv)}
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !can(__d, r) {
			out.Fail(w, http.StatusForbidden, "no plan")
			return
		}
		_id := auth.ID(r)
		_v := View{Plugs: []Plug{}, Gh: __d.Cfg.Gh}
		_d, _err := __d.Devs.Get(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if _d != nil {
			_v.Nick, _v.Bio = _d.Nick, _d.Bio
		}
		_ps, _err := __d.Devs.Plugs(r.Context(), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_slugs := make([]string, 0, len(_ps))
		for _, _p := range _ps {
			_slugs = append(_slugs, _p.Slug)
		}
		_n, _ := __d.Mods.Users(r.Context(), _slugs)
		for _, _p := range _ps {
			_v.Plugs = append(_v.Plugs, view(_p, _n[_p.Slug]))
		}
		out.JSON(w, _v)
	}
}

func Post(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !can(__d, r) {
			out.Fail(w, http.StatusForbidden, "no plan")
			return
		}
		var _in In
		if !out.Read(w, r, &_in) {
			return
		}
		_in.Nick = strings.TrimLeft(strings.TrimSpace(_in.Nick), "@")
		_in.Bio = strings.TrimSpace(_in.Bio)
		if !__nick__.MatchString(_in.Nick) {
			out.Fail(w, http.StatusBadRequest, "bad nick")
			return
		}
		if utf8.RuneCountInString(_in.Bio) > 120 {
			out.Fail(w, http.StatusBadRequest, "bio too long")
			return
		}
		_err := __d.Devs.Set(r.Context(), devs.Dev{User: auth.ID(r), Nick: _in.Nick, Bio: _in.Bio})
		if errors.Is(_err, devs.ErrTaken) {
			out.Fail(w, http.StatusConflict, "nick taken")
			return
		}
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		out.JSON(w, _in)
	}
}

func mine(__d *deps.D, w http.ResponseWriter, r *http.Request) *devs.Plug {
	if !can(__d, r) {
		out.Fail(w, http.StatusForbidden, "no plan")
		return nil
	}
	_id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_p, _err := __d.Devs.Plug(r.Context(), _id)
	if _err != nil {
		out.Oops(w, _err)
		return nil
	}
	if _p == nil || _p.User != auth.ID(r) {
		out.Fail(w, http.StatusNotFound, "not found")
		return nil
	}
	return _p
}

func Show(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_p := mine(__d, w, r)
		if _p == nil {
			return
		}
		_n, _ := __d.Mods.Users(r.Context(), []string{_p.Slug})
		_rs, _err := __d.Devs.History(r.Context(), _p.ID)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_o := One{Plug: view(*_p, _n[_p.Slug]), Reqs: make([]Req, 0, len(_rs))}
		for _, _q := range _rs {
			_o.Reqs = append(_o.Reqs, Req{ID: _q.ID, Kind: _q.Kind, Name: _q.Name, About: _q.About, Ver: _q.Ver, Note: _q.Note, Repo: _q.Repo, Pic: _q.Pic,
				Status: _q.Status, Reason: _q.Reason, Created: _q.Created.Unix(), Done: _q.Done.Unix()})
		}
		out.JSON(w, _o)
	}
}

func read(w http.ResponseWriter, r *http.Request) (*PlugIn, []byte, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, __pic__*2)
	var _in PlugIn
	if !out.Read(w, r, &_in) {
		return nil, nil, false
	}
	var _ok [5]bool
	_in.Name, _ok[0] = short(_in.Name, 32)
	_in.About, _ok[1] = short(_in.About, 300)
	_in.Ver, _ok[2] = short(strings.ReplaceAll(_in.Ver, ",", "."), 16)
	_in.Note, _ok[3] = short(_in.Note, 1000)
	_in.Slug = strings.ToLower(strings.TrimSpace(_in.Slug))
	_ok[4] = utf8.RuneCountInString(_in.Name) >= 2
	for _, _x := range _ok {
		if !_x {
			out.Fail(w, http.StatusBadRequest, "bad fields")
			return nil, nil, false
		}
	}
	_pic, _good := Pic(_in.Pic)
	if !_good {
		out.Fail(w, http.StatusBadRequest, "bad pic")
		return nil, nil, false
	}
	return &_in, _pic, true
}

func limit(__d *deps.D, w http.ResponseWriter, r *http.Request) (*devs.Dev, bool) {
	_d, _err := __d.Devs.Get(r.Context(), auth.ID(r))
	if _err != nil {
		out.Oops(w, _err)
		return nil, false
	}
	if _d == nil {
		out.Fail(w, http.StatusConflict, "no nick")
		return nil, false
	}
	if _n, _ := __d.Devs.Waiting(r.Context(), _d.User); _n >= 5 {
		out.Fail(w, http.StatusTooManyRequests, "too many")
		return nil, false
	}
	return _d, true
}

func New(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !can(__d, r) {
			out.Fail(w, http.StatusForbidden, "no plan")
			return
		}
		_d, _ok := limit(__d, w, r)
		if !_ok {
			return
		}
		_in, _pic, _ok := read(w, r)
		if !_ok {
			return
		}
		if !__slug__.MatchString(_in.Slug) {
			out.Fail(w, http.StatusBadRequest, "bad slug")
			return
		}
		_repo, _ok := Repo(_in.Repo)
		if !_ok {
			out.Fail(w, http.StatusBadRequest, "bad repo")
			return
		}
		_q := devs.Req{User: _d.User, Name: _in.Name, About: _in.About, Ver: _in.Ver, Note: _in.Note, Repo: _repo, Pic: _pic != nil}
		_pid, _qid, _err := __d.Devs.Create(r.Context(), devs.Plug{User: _d.User, Slug: _in.Slug, Name: _in.Name, About: _in.About, Repo: _repo, Ver: _in.Ver}, _q, _pic)
		if errors.Is(_err, devs.ErrSlug) {
			out.Fail(w, http.StatusConflict, "slug taken")
			return
		}
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_q.ID, _q.Plug, _q.Kind = _qid, _pid, "new"
		review.Ask(r.Context(), __d, _d, _q, _in.Slug)
		out.JSON(w, map[string]int64{"id": _pid})
	}
}

func Change(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_p := mine(__d, w, r)
		if _p == nil {
			return
		}
		_d, _ok := limit(__d, w, r)
		if !_ok {
			return
		}
		if _busy, _ := __d.Devs.Busy(r.Context(), _p.ID); _busy {
			out.Fail(w, http.StatusConflict, "busy")
			return
		}
		_in, _pic, _ok := read(w, r)
		if !_ok {
			return
		}
		_q := devs.Req{User: _d.User, Plug: _p.ID, Name: _in.Name, About: _in.About, Note: _in.Note, Pic: _pic != nil}
		switch _in.Kind {
		case "edit":
			_q.Kind, _q.Ver, _q.Repo = "edit", _p.Ver, _p.Repo
		case "update":
			_repo, _good := Repo(_in.Repo)
			if !_good {
				out.Fail(w, http.StatusBadRequest, "bad repo")
				return
			}
			if _in.Ver == "" || (_in.Ver == _p.Ver && _p.Live) {
				out.Fail(w, http.StatusBadRequest, "same ver")
				return
			}
			_q.Kind, _q.Ver, _q.Repo = "update", _in.Ver, _repo
		default:
			out.Fail(w, http.StatusBadRequest, "bad kind")
			return
		}
		_id, _err := __d.Devs.Add(r.Context(), _q, _pic)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_q.ID = _id
		review.Ask(r.Context(), __d, _d, _q, _p.Slug)
		out.JSON(w, map[string]int64{"id": _id})
	}
}

func ReqPic(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		_b, _u, _err := __d.Devs.ReqPic(r.Context(), _id)
		_me := auth.ID(r)
		if _err != nil || _b == nil || (_me != _u && _me != __d.Cfg.Own) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", http.DetectContentType(_b))
		w.Header().Set("Cache-Control", "private, max-age=86400")
		_, _ = w.Write(_b)
	}
}
