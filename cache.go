package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// Cacher that adapter needs to implement
type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, data string, expire time.Duration) error
	Delete(ctx context.Context, key string) error
}

type Options struct {
	expire time.Duration
}

type Option func(*Options)

// WithExpire sets the expire time for the cache
func WithExpire(expire time.Duration) Option {
	return func(o *Options) {
		o.expire = expire
	}
}

// Cache caches the data which is returned by fn with the given key
func Cache[T any](ctx context.Context, c Cacher, key string, fn func() (T, error), opts ...Option) (T, error) {
	var data T
	options := &Options{
		expire: 30 * time.Minute,
	}

	for _, opt := range opts {
		opt(options)
	}

	if info, err := c.Get(ctx, key); err == nil {
		return unmarshal[T](info)
	}

	data, err := fn()
	if err != nil {
		return data, err
	}

	val, err := marshal(data)
	if err != nil {
		return data, err
	}

	if err := c.Set(ctx, key, val, options.expire); err != nil {
		log.Println("[cache] set error:", err)
	}
	return data, nil
}

func marshal(data any) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshal[T any](data string) (T, error) {
	var m T
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return m, err
	}
	return m, nil
}
