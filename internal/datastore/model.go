package datastore

import "time"

const (
	Job   ItemType = "job"
	Story ItemType = "story"
)

type ItemType string

type Item struct {
	ID        int
	Type      ItemType
	Title     string
	Content   string
	URL       string
	Score     int
	CreatedBy string
	CreatedAt time.Time
}
