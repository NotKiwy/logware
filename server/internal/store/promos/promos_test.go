package promos

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"logware/internal/pay"
	"logware/internal/store/db"
)

func TestUse(t *testing.T) {
	_url := os.Getenv("LW_TEST_DB")
	if _url == "" {
		t.Skip()
	}
	_ctx := context.Background()
	_p, _err := db.Open(_ctx, _url)
	if _err != nil {
		t.Fatal(_err)
	}
	_, _ = _p.Exec(_ctx, `delete from promo_uses; delete from promos; delete from subs where user_id in (1, 2, 3)`)
	_, _ = _p.Exec(_ctx, `insert into promos (code, plan, days, cap) values ('WEEK', 'pro', 7, 2), ('DEV', 'dev', 3, 0)`)
	_, _ = _p.Exec(_ctx, `insert into promos (code, plan, days, ends) values ('OLD', 'pro', 7, now() - interval '1 day')`)
	_r := New(_p)
	_pl, _u, _err := _r.Use(_ctx, " week ", 1, pay.Merge)
	if _err != nil || _pl != "pro" || time.Until(_u) < 6*24*time.Hour {
		t.Fatal("first", _pl, _u, _err)
	}
	if _, _, _err = _r.Use(_ctx, "WEEK", 1, pay.Merge); !errors.Is(_err, ErrUsed) {
		t.Fatal("reuse", _err)
	}
	if _, _, _err = _r.Use(_ctx, "WEEK", 2, pay.Merge); _err != nil {
		t.Fatal("second user", _err)
	}
	if _, _, _err = _r.Use(_ctx, "WEEK", 3, pay.Merge); !errors.Is(_err, ErrOver) {
		t.Fatal("cap", _err)
	}
	if _, _, _err = _r.Use(_ctx, "OLD", 3, pay.Merge); !errors.Is(_err, ErrOver) {
		t.Fatal("ends", _err)
	}
	if _, _, _err = _r.Use(_ctx, "NOPE", 3, pay.Merge); !errors.Is(_err, ErrNone) {
		t.Fatal("none", _err)
	}
	_pl, _u2, _err := _r.Use(_ctx, "dev", 1, pay.Merge)
	if _err != nil || _pl != "dev" || time.Until(_u2) > 7*24*time.Hour || time.Until(_u2) < 6*24*time.Hour {
		t.Fatal("upgrade stacks", _pl, _u2, _err)
	}
	_, _ = _p.Exec(_ctx, `update promos set off = true where code = 'DEV'`)
	if _, _, _err = _r.Use(_ctx, "DEV", 2, pay.Merge); !errors.Is(_err, ErrNone) {
		t.Fatal("off", _err)
	}
	var _ch string
	_ = _p.QueryRow(_ctx, `select charge from subs where user_id = 1`).Scan(&_ch)
	if _ch != "promo:DEV" || pay.Star(_ch) {
		t.Fatal("charge", _ch)
	}
}
