package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/init.sql
var __schema__ string

func Open(__ctx context.Context, __url string) (*pgxpool.Pool, error) {
	_p, _err := pgxpool.New(__ctx, __url)
	if _err != nil {
		return nil, _err
	}
	if _err = _p.Ping(__ctx); _err != nil {
		_p.Close()
		return nil, _err
	}
	if _, _err = _p.Exec(__ctx, __schema__); _err != nil {
		_p.Close()
		return nil, _err
	}
	return _p, nil
}
