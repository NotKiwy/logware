package pay

import (
	"strconv"
	"strings"
	"time"
)

const Period = 30 * 24 * 60 * 60

type Plan struct {
	ID    string
	Name  string
	Price int
	Rub   int
	Usd   string
}

var __plans__ = map[string]Plan{
	"pro": {ID: "pro", Name: "Logware Pro", Price: 149, Rub: 149, Usd: "1.79"},
	"dev": {ID: "dev", Name: "Logware Developer", Price: 299, Rub: 299, Usd: "3.59"},
}

func Get(__id string) (Plan, bool) {
	_p, _ok := __plans__[__id]
	return _p, _ok
}

func Rank(__id string) int {
	switch __id {
	case "pro":
		return 1
	case "dev":
		return 2
	}
	return 0
}

func Payload(__plan string, __user int64) string {
	return __plan + ":" + strconv.FormatInt(__user, 10)
}

func Tagged(__plan string, __user int64, __tag string) string {
	return Payload(__plan, __user) + "@" + __tag
}

func Tag(__s string) (string, string) {
	if _i := strings.LastIndex(__s, "@"); _i >= 0 {
		return __s[:_i], __s[_i+1:]
	}
	return __s, ""
}

func Parse(__s string) (Plan, int64, bool) {
	__s, _ = Tag(__s)
	_id, _u, _ok := strings.Cut(__s, ":")
	if !_ok {
		return Plan{}, 0, false
	}
	_p, _ok := Get(_id)
	_n, _err := strconv.ParseInt(_u, 10, 64)
	return _p, _n, _ok && _err == nil
}

func Star(__charge string) bool {
	return __charge != "" && !strings.HasPrefix(__charge, "admin") && !strings.Contains(__charge, ":")
}

func Merge(__plan string, __until time.Time, __new string, __add time.Duration, __now time.Time) (string, time.Time) {
	_left := __until.Sub(__now)
	if _left <= 0 || __plan == "" {
		return __new, __now.Add(__add)
	}
	if __plan == __new {
		return __new, __until.Add(__add)
	}
	_a, _ := Get(__plan)
	_b, _ := Get(__new)
	if _a.Price == 0 || _b.Price == 0 {
		return __new, __now.Add(__add)
	}
	if Rank(__new) > Rank(__plan) {
		return __new, __now.Add(__add + time.Duration(float64(_left)*float64(_a.Price)/float64(_b.Price)))
	}
	return __plan, __until.Add(time.Duration(float64(__add) * float64(_b.Price) / float64(_a.Price)))
}
