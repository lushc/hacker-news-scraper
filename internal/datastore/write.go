package datastore

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	//go:embed queries/insert.sql
	insertQuery string
)

type Writer interface {
	Insert(ctx context.Context, item Item) error
	Close()
}

type DBWriter struct {
	pool *pgxpool.Pool
}

func NewDBWriter(ctx context.Context, maxConns int) (*DBWriter, error) {
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

	return &DBWriter{pool: pool}, nil
}

func (d DBWriter) Insert(ctx context.Context, item Item) error {
	if _, err := d.pool.Exec(ctx, insertQuery, item.ID, item.Type, item.Title, item.Content, item.URL, item.Score, item.CreatedBy, item.CreatedAt); err != nil {
		return fmt.Errorf("insert item %d: %w", item.ID, err)
	}

	return nil
}

func (d DBWriter) Close() {
	d.pool.Close()
}
