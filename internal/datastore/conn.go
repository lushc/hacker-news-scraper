package datastore

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBConn struct {
	*pgxpool.Pool
}

func NewDBConn(ctx context.Context, maxConns int) (*DBConn, error) {
	// TODO: this + maxConns could be viper config based instead
	connUrl := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxConns)
	config.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &DBConn{pool}, nil
}
