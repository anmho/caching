package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache[T any] struct {
	redisClient *redis.Client
}

func New[T any](
	redisClient *redis.Client,
) *Cache[T] {
	return &Cache[T]{
		redisClient: redisClient,
	}
}

func (c *Cache[T]) InvalidateKey(ctx context.Context, key string) error {

	cmd := c.redisClient.Del(ctx, key)

	err := cmd.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache[T]) WriteItem(ctx context.Context, key string, data *T) error {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cmd := c.redisClient.Set(ctx, key, b, time.Minute*0)
	err = cmd.Err()

	if err != nil {
		return err
	}

	return nil
}

type ReadCacheResult[T any] struct {
	Data     *T
	CacheHit bool
}

func (c *Cache[T]) ReadItem(ctx context.Context, key string) (ReadCacheResult[T], error) {
	cmd := c.redisClient.Get(ctx, key)
	err := cmd.Err()
	if errors.Is(err, redis.Nil) {
		return ReadCacheResult[T]{
			Data:     nil,
			CacheHit: false,
		}, nil
	}

	b, err := cmd.Bytes()
	if err != nil {
		return ReadCacheResult[T]{
			CacheHit: false,
		}, err
	}

	data := new(T)

	err = json.Unmarshal(b, data)
	if err != nil {
		return ReadCacheResult[T]{
			CacheHit: true,
		}, err
	}

	return ReadCacheResult[T]{
		Data:     data,
		CacheHit: true,
	}, nil
}
