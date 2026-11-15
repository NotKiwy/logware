package cfg

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Cfg struct {
	Tok  string
	DB   string
	Pub  string
	Addr string
	Web  string
	Own  int64
	Gh   string
	Cb   string
	Xr   string
	XrWh string
	Tag  string
	SupT string
}

func Load() (*Cfg, error) {
	Env(".env", "../.env")
	_c := &Cfg{
		Tok:  os.Getenv("BOT_TOKEN"),
		DB:   or(os.Getenv("DB_URL"), "postgres://localhost:5432/logware?sslmode=disable"),
		Pub:  strings.TrimRight(os.Getenv("PUBLIC_URL"), "/"),
		Addr: or(os.Getenv("ADDR"), ":8080"),
		Web:  or(os.Getenv("WEB_DIR"), "../web/dist"),
	}
	_c.Own, _ = strconv.ParseInt(os.Getenv("OWNER"), 10, 64)
	_c.Gh = strings.TrimPrefix(os.Getenv("GITHUB"), "@")
	_c.Cb = os.Getenv("CRYPTOBOT_TOKEN")
	_c.Xr = os.Getenv("XROCKET_TOKEN")
	_c.XrWh = os.Getenv("XROCKET_SECRET")
	_c.SupT = os.Getenv("SUPPORT_TOKEN")
	_c.Tag, _, _ = strings.Cut(strings.TrimPrefix(strings.TrimPrefix(_c.Pub, "https://"), "http://"), ".")
	if _c.Tok == "" {
		return nil, errors.New("BOT_TOKEN is empty, put it into .env")
	}
	return _c, nil
}

func or(__v, __def string) string {
	if __v == "" {
		return __def
	}
	return __v
}
