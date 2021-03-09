package datastore

import (
	"context"
	_ "embed"

	"github.com/georgysavva/scany/pgxscan"
)

var (
	//go:embed queries/all.sql
	allQuery string
	//go:embed queries/by_type.sql
	byTypeQuery string
)

// Reader is an interface for querying for items in the datastore
type Reader interface {
	All(ctx context.Context) ([]*Item, error)
	ByItemType(ctx context.Context, itemType ItemType) ([]*Item, error)
	Close()
}

// DBReader is a pgsql database reader
type DBReader struct {
	conn DBConn
}

// NewDBReader creates a new database reader
func NewDBReader(conn DBConn) *DBReader {
	return &DBReader{conn: conn}
}

// All returns all top story items from the DB
func (d DBReader) All(ctx context.Context) (items []*Item, err error) {
	err = pgxscan.Select(ctx, d.conn, &items, allQuery)
	return items, err
}

// ByItemType returns all top story items with the given item type from the DB
func (d DBReader) ByItemType(ctx context.Context, itemType ItemType) (items []*Item, err error) {
	err = pgxscan.Select(ctx, d.conn, &items, byTypeQuery, itemType)
	return items, err
}

// Close closes the underlying DB connection
func (d DBReader) Close() {
	d.conn.Close()
}
