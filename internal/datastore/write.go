package datastore

import (
	"context"
	_ "embed"
	"fmt"
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
	conn DBConn
}

func NewDBWriter(conn DBConn) *DBWriter {
	return &DBWriter{conn: conn}
}

func (d DBWriter) Insert(ctx context.Context, item Item) error {
	if _, err := d.conn.Exec(ctx, insertQuery, item.ID, item.Type, item.Title, item.Content, item.URL, item.Score, item.CreatedBy, item.CreatedAt); err != nil {
		return fmt.Errorf("insert item %d: %w", item.ID, err)
	}

	return nil
}

func (d DBWriter) Close() {
	d.conn.Close()
}
