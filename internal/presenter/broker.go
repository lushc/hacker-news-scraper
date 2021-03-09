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

// EventStream is a channel of SSE messages
type EventStream chan string

// Source is a function which sends items to a channel
type Source func(items chan<- datastore.Item, errs chan<- error, ephemeral bool)

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
	// channel where errors are pushed to when fetching items
	errs chan error
}

// NewBroker creates a new broker for the event source
func NewBroker(source Source) *Broker {
	return &Broker{
		source:      source,
		logger:      logrus.New(),
		subscribers: make(map[EventStream]interface{}),
		newSubs:     make(chan EventStream),
		expiredSubs: make(chan EventStream),
		items:       make(chan datastore.Item),
		errs:        make(chan error),
	}
}

// Handle will start the broker and return a handler for echo to invoke when an endpoint is hit
func (b *Broker) Handle(ctx context.Context) echo.HandlerFunc {
	go b.run(ctx)
	return b.handler
}

// Refresh will source items to be broadcast to subscribers
func (b *Broker) Refresh() {
	b.source(b.items, b.errs, false)
}

// run will manage subscribers and broadcast events to them
func (b *Broker) run(ctx context.Context) {
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
				b.logger.Error(fmt.Errorf("encode event for subscribers: %w", err))
				break
			}

			for ch := range b.subscribers {
				ch <- evt
			}
		case err := <-b.errs:
			b.logger.Error(fmt.Errorf("error from broker channel: %w", err))
		}
	}
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
		items := make(chan datastore.Item)
		errs := make(chan error)
		go b.source(items, errs, true)

		for {
			select {
			case <-c.Request().Context().Done():
				return
			case item, ok := <-items:
				if !ok {
					return
				}

				evt, err := encode(item)
				if err != nil {
					b.logger.Error(fmt.Errorf("encode event for client: %w", err))
					break
				}

				subscriber <- evt
			case err, ok := <-errs:
				if !ok {
					return
				}

				b.logger.Error(fmt.Errorf("error from handler channel: %w", err))
			}
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
				return fmt.Errorf("write SSE response: %w", err)
			}

			res.Flush()
		}
	}
}

// encode will create an SSE message for an item
func encode(item datastore.Item) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal item to JSON: %w", err)
	}

	return fmt.Sprintf("id: %d\ndata: %s\n\n", item.ID, data), nil
}
