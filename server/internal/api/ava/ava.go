package ava

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
)

const (
	__ttl__  = 6 * time.Hour
	__miss__ = time.Hour
	__max__  = 3000
)

type _Pic struct {
	data []byte
	till time.Time
}

var (
	__mu   sync.Mutex
	__pics = map[int64]_Pic{}
)

func load(__d *deps.D, r *http.Request, __id int64) []byte {
	__mu.Lock()
	_p, _ok := __pics[__id]
	__mu.Unlock()
	if _ok && time.Now().Before(_p.till) {
		return _p.data
	}
	var _data []byte
	_f, _err := __d.Tg.Photo(r.Context(), __id)
	if _err == nil && _f != nil && _f.Path != "" {
		_data, _ = __d.Tg.Fetch(r.Context(), _f.Path)
	}
	_ttl := __ttl__
	if len(_data) == 0 {
		_ttl = __miss__
	}
	__mu.Lock()
	if len(__pics) >= __max__ {
		__pics = map[int64]_Pic{}
	}
	__pics[__id] = _Pic{data: _data, till: time.Now().Add(_ttl)}
	__mu.Unlock()
	return _data
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id, _err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if _err != nil || _id <= 0 {
			out.Fail(w, http.StatusBadRequest, "bad id")
			return
		}
		_me := auth.ID(r)
		if _id != _me {
			_ok, _err := __d.Msgs.HasChat(r.Context(), _me, _id)
			if _err != nil {
				out.Oops(w, _err)
				return
			}
			if !_ok {
				out.Fail(w, http.StatusNotFound, "not found")
				return
			}
		}
		Send(__d, w, r, _id)
	}
}

func Send(__d *deps.D, w http.ResponseWriter, r *http.Request, __id int64) {
	_data := load(__d, r, __id)
	if len(_data) == 0 {
		out.Fail(w, http.StatusNotFound, "no photo")
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	_, _ = w.Write(_data)
}
