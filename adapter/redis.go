package adapter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kimistar/go-cache"
)

var _ cache.Cacher = (*Redis)(nil)

type Redis struct {
	client *redis.Client
}

// NewRedis creates redis cache adapter
func NewRedis(client *redis.Client) cache.Cacher {
	return &Redis{
		client: client,
	}
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, data string, expiration time.Duration) error {
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
