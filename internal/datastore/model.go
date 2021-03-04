package datastore

import "time"

type Item struct {
	ID        int
	Type      string
	Title     string
	Content   string
	URL       string
	Score     int
	CreatedBy string
	CreatedAt time.Time
}
