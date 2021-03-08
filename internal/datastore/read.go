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

type Reader interface {
	All(ctx context.Context) ([]*Item, error)
	ByItemType(ctx context.Context, itemType ItemType) ([]*Item, error)
	Close()
}

type DBReader struct {
	conn DBConn
}

func NewDBReader(conn DBConn) *DBReader {
	return &DBReader{conn: conn}
}

func (d DBReader) All(ctx context.Context) (items []*Item, err error) {
	err = pgxscan.Select(ctx, d.conn, &items, allQuery)
	return items, err
}

func (d DBReader) ByItemType(ctx context.Context, itemType ItemType) (items []*Item, err error) {
	err = pgxscan.Select(ctx, d.conn, &items, byTypeQuery, itemType)
	return items, err
}

func (d DBReader) Close() {
	d.conn.Close()
}
