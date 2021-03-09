package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
	"github.com/sirupsen/logrus"
)

// Worker is a worker for fetching a Hacker News item and persisting it to the datastore
type Worker struct {
	client Client
	writer datastore.Writer
	logger *logrus.Logger
}

// NewWorker creates a new worker
func NewWorker(client Client, writer datastore.Writer) *Worker {
	return &Worker{
		client: client,
		writer: writer,
		logger: logrus.New(),
	}
}

// Run is meant to be ran concurrently, receiving item IDs from a channel to fetch and persist
func (w Worker) Run(ctx context.Context, items <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-items:
			if !ok {
				return
			}

			item, err := w.client.Item(ctx, id)
			if err != nil {
				w.logger.Error(err)
				break
			}

			if item.Deleted || item.Dead {
				break
			}

			record := datastore.Item{
				ID:        item.ID,
				Type:      datastore.ItemType(item.Type),
				Title:     item.Title,
				Content:   item.Text,
				URL:       item.URL,
				Score:     item.Score,
				CreatedBy: item.By,
				CreatedAt: time.Unix(item.Time, 0),
			}

			if err := w.writer.Insert(ctx, record); err != nil {
				w.logger.Error(err)
				break
			}
		}
	}
}
