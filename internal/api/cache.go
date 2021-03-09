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

type CachedReader struct {
	cache  *cache.Cache
	reader datastore.Reader
}

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

func (c CachedReader) Close() {
	c.reader.Close()
}
