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

// Writer is an interface for persisting items to the datastore
type Writer interface {
	Insert(ctx context.Context, item Item) error
	Close()
}

// DBWriter is a pgsql database writer
type DBWriter struct {
	conn DBConn
}

// NewDBWriter creates a new database writer
func NewDBWriter(conn DBConn) *DBWriter {
	return &DBWriter{conn: conn}
}

// Insert will insert the given item into the DB
func (d DBWriter) Insert(ctx context.Context, item Item) error {
	if _, err := d.conn.Exec(ctx, insertQuery, item.ID, item.Type, item.Title, item.Content, item.URL, item.Score, item.CreatedBy, item.CreatedAt); err != nil {
		return fmt.Errorf("insert item %d: %w", item.ID, err)
	}

	return nil
}

// Close closes the underlying database connection
func (d DBWriter) Close() {
	d.conn.Close()
}
