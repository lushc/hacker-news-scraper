package presenter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
)

// EventStream is a channel of SSE events
type EventStream chan string

// Source is a function which provides items
type Source func() ([]datastore.Item, error)

// Broker is responsible for managing client connections and sending events
type Broker struct {
	source Source
	logger *logrus.Logger
	// map of subscriber channels which events can be pushed to
	subscribers map[EventStream]interface{}
	// channel which new subscribers will be pushed to
	newSubs chan EventStream
	// channel which disconnected subscribers will be pushed to
	expiredSubs chan EventStream
	// channel where items are pushed to for broadcasting to subscribers as an event
	items chan datastore.Item
}

func NewBroker(source Source) *Broker {
	return &Broker{
		source:      source,
		logger:      logrus.New(),
		subscribers: make(map[EventStream]interface{}),
		newSubs:     make(chan EventStream),
		expiredSubs: make(chan EventStream),
		items:       make(chan datastore.Item),
	}
}

func (b *Broker) Handle(ctx context.Context) echo.HandlerFunc {
	b.run(ctx)
	return b.handler
}

// run will begin a goroutine that manages subscribers and broadcasting events to them
func (b *Broker) run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch := <-b.newSubs:
				// TODO: mux needed?
				b.subscribers[ch] = nil
			case ch := <-b.expiredSubs:
				// TODO: mux needed?
				delete(b.subscribers, ch)
				close(ch)
			case item := <-b.items:
				evt, err := encode(item)
				if err != nil {
					b.logger.Error(err)
					break
				}

				for ch := range b.subscribers {
					ch <- evt
				}
			}
		}
	}()
}

// handler will stream event data to the client of the invoking HTTP request
func (b *Broker) handler(c echo.Context) error {
	// add a new channel for the connecting client
	subscriber := make(EventStream)
	b.newSubs <- subscriber

	// send SSE headers
	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.WriteHeader(http.StatusOK)

	// start pushing items to the client
	go func() {
		items, err := b.source()
		if err != nil {
			b.logger.Error(fmt.Errorf("failed to fetch items for client: %w", err))
			return
		}

		for _, item := range items {
			evt, err := encode(item)
			if err != nil {
				b.logger.Error(fmt.Errorf("failed to encode event for client: %w", err))
				return
			}

			subscriber <- evt
		}
	}()

	// block and stream event responses while the connection remains open
	for {
		select {
		case <-c.Request().Context().Done():
			b.expiredSubs <- subscriber
			return nil
		case evt, ok := <-subscriber:
			if !ok {
				break
			}

			if _, err := res.Write([]byte(evt)); err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}

			res.Flush()
		}
	}
}

// encode will create an SSE string for an item
func encode(item datastore.Item) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("failed to marshal item to JSON: %w", err)
	}

	return fmt.Sprintf("id: %d\ndata: %s\n\n", item.ID, data), nil
}
