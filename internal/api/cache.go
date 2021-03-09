package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
)

const (
	redisHostEnv = "REDIS_HOST"
	ttl          = 5 * time.Minute
)

var (
	errRedisHostEnv = fmt.Errorf("missing env var %s", redisHostEnv)
)

// CachedReader is a wrapper to add Redis-backed caching to an underlying datastore.Reader
type CachedReader struct {
	cache  *cache.Cache
	reader datastore.Reader
}

// NewCachedReader creates a new CachedReader with a pre-configured Redis client
func NewCachedReader(reader datastore.Reader) (*CachedReader, error) {
	// TODO: viper config instead
	host, ok := os.LookupEnv(redisHostEnv)
	if !ok {
		return nil, errRedisHostEnv
	}

	redisCache := cache.New(&cache.Options{
		Redis: redis.NewClient(&redis.Options{
			Addr: host,
		}),
	})

	return &CachedReader{
		cache:  redisCache,
		reader: reader,
	}, nil
}

// All caches and returns all items from the reader
func (c CachedReader) All(ctx context.Context) (items []*datastore.Item, err error) {
	err = c.cache.Once(&cache.Item{
		Key:   "all-stories",
		Value: &items,
		TTL:   ttl,
		Do: func(item *cache.Item) (interface{}, error) {
			return c.reader.All(ctx)
		},
	})
	return items, err
}

// ByItemType caches and returns items of the given type from the reader
func (c CachedReader) ByItemType(ctx context.Context, itemType datastore.ItemType) (items []*datastore.Item, err error) {
	err = c.cache.Once(&cache.Item{
		Key:   fmt.Sprintf("%s-stories", itemType),
		Value: &items,
		TTL:   ttl,
		Do: func(item *cache.Item) (interface{}, error) {
			return c.reader.ByItemType(ctx, itemType)
		},
	})
	return items, err
}

// Close closes the underlying reader connection
func (c CachedReader) Close() {
	c.reader.Close()
}
