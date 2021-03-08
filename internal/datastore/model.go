package datastore

import (
	"time"

	pb "github.com/lushc/hacker-news-scraper/protobufs"
)

const (
	Job   ItemType = "job"
	Story ItemType = "story"
)

var (
	EnumTypes = map[pb.Type]ItemType{
		pb.Type_JOB:   Job,
		pb.Type_STORY: Story,
	}
)

type ItemType string

// TODO: refactor to just use the protobuf type?
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

func Itop(item Item) *pb.Item {
	var itemType pb.Type
	for k, v := range EnumTypes {
		if v == item.Type {
			itemType = k
			break
		}
	}

	return &pb.Item{
		Id:        int32(item.ID),
		Type:      itemType,
		Title:     item.Title,
		Content:   item.Content,
		Url:       item.URL,
		Score:     int32(item.Score),
		CreatedBy: item.CreatedBy,
		CreatedAt: item.CreatedAt.Unix(),
	}
}

func Ptoi(proto *pb.Item) *Item {
	return &Item{
		ID:        int(proto.Id),
		Type:      EnumTypes[proto.Type],
		Title:     proto.Title,
		Content:   proto.Content,
		URL:       proto.Url,
		Score:     int(proto.Score),
		CreatedBy: proto.CreatedBy,
		CreatedAt: time.Unix(proto.CreatedAt, 0),
	}
}
