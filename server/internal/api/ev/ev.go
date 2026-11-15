package ev

import (
	"net/http"
	"strconv"

	"logware/internal/api/auth"
	"logware/internal/api/out"
	"logware/internal/deps"
)

type Item struct {
	ID    int64  `json:"id"`
	Ver   int    `json:"ver"`
	Out   bool   `json:"out"`
	Kind  string `json:"kind"`
	Body  []byte `json:"body"`
	Media bool   `json:"media"`
	At    int64  `json:"at"`
}

type Ev struct {
	Kind  string `json:"kind"`
	At    int64  `json:"at"`
	Total int    `json:"total"`
	More  bool   `json:"more"`
	Items []Item `json:"items"`
}

func num(__r *http.Request, __k string, __def, __max int) int {
	_n, _err := strconv.Atoi(__r.URL.Query().Get(__k))
	if _err != nil || _n < 0 {
		return __def
	}
	return min(_n, __max)
}

func Get(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := auth.ID(r)
		_e, _err := __d.Evs.Get(r.Context(), r.PathValue("id"), _id)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if _e == nil {
			out.Fail(w, http.StatusNotFound, "not found")
			return
		}
		_lim := num(r, "n", 50, 200)
		_rows, _err := __d.Msgs.Page(r.Context(), _id, _e.Refs, num(r, "from", 0, 1<<30), _lim+1)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		_more := len(_rows) > _lim
		if _more {
			_rows = _rows[:_lim]
		}
		_out := Ev{Kind: _e.Kind, At: _e.Created.Unix(), Total: len(_e.Refs), More: _more, Items: make([]Item, 0, len(_rows))}
		for _, _r := range _rows {
			_out.Items = append(_out.Items, Item{ID: _r.ID, Ver: _r.Ver, Out: _r.Out, Kind: _r.Kind, Body: _r.Body, Media: _r.Has, At: _r.Created.Unix()})
		}
		out.JSON(w, _out)
	}
}

func Media(__d *deps.D) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_n, _err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if _err != nil {
			out.Fail(w, http.StatusBadRequest, "bad id")
			return
		}
		_b, _err := __d.Msgs.Media(r.Context(), auth.ID(r), _n)
		if _err != nil {
			out.Oops(w, _err)
			return
		}
		if len(_b) == 0 {
			out.Fail(w, http.StatusNotFound, "not found")
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(_b)
	}
}
