package clock

import (
	"context"
	"errors"

	"logware/internal/deps"
	"logware/internal/mods/mod"
)

const ID = "clock"

var Mod = &mod.Mod{
	ID:     ID,
	Name:   "Время в имени",
	Desc:   "Часы рядом с твоим именем, обновляются каждую минуту",
	Icon:   "clock",
	Color:  "#5856d6",
	Author: "Logware",
	Hidden: true,
	Perms: []mod.Perm{
		{Right: "name", Why: "Чтобы менять твоё имя каждую минуту"},
	},
	Opts: []mod.Opt{
		{Key: "tpl", Name: "Шаблон", Def: "[%H:%M]", Text: true, Max: 40},
		{Key: "pos", Name: "Где показывать", Def: "after", Choices: []mod.Choice{
			{V: "before", Label: "Перед именем"},
			{V: "after", Label: "После имени"},
		}},
		{Key: "first", Name: "Имя", Text: true, Max: 64},
		{Key: "last", Name: "Фамилия", Text: true, Max: 64},
		{Key: "tz", Name: "Часовой пояс", Def: "Europe/Moscow", Text: true, Max: 64, Hide: true},
	},
	Ready: func(*deps.D) bool { return true },
	Set:   set,
}

var ErrConn = errors.New("Сначала подключи бота: Профиль > Изменить > ИИ-автоматизация")

func set(__ctx context.Context, __d *deps.D, __user int64, __on bool, __c mod.Cfg) (mod.Cfg, error) {
	_c, _err := __d.Conns.ByUser(__ctx, __user)
	if _err != nil {
		return nil, _err
	}
	if _c == nil {
		if __on {
			return nil, ErrConn
		}
		return __c, nil
	}
	if !__on {
		forget(__user)
		if __c["first"] != "" {
			return __c, __d.Tg.BizName(__ctx, _c.ID, __c["first"], __c["last"])
		}
		return __c, nil
	}
	_tc, _err := __d.Tg.Conn(__ctx, _c.ID)
	if _err != nil {
		return nil, _err
	}
	if _err = __d.Conns.Put(__ctx, _tc); _err != nil {
		return nil, _err
	}
	if !_tc.On {
		return nil, ErrConn
	}
	if !_tc.Rights.Name {
		return nil, errors.New("Разреши боту менять имя: Профиль > Изменить > ИИ-автоматизация > @" + __d.Nick)
	}
	if __c["first"] == "" {
		__c["first"], __c["last"] = _tc.User.First, _tc.User.Last
	}
	_f, _l := Name(__c, now(__c))
	forget(__user)
	return __c, __d.Tg.BizName(__ctx, _c.ID, _f, _l)
}
