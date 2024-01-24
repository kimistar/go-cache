package adapter

import (
	"context"
	"errors"
	"log"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/kimistar/go-cache"
)

var ErrNoData = errors.New("no data")
var ErrExpire = errors.New("data is expired")

// NewLocal creates local cache adapter, cache type is string
func NewLocal(opts ...LocalOption) cache.Cacher {
	return NewLocalCache[string](opts...)
}

// NewLocalDefault mainly used for wire injection without additional arguments
func NewLocalDefault() cache.Cacher {
	return NewLocalCache[string]()
}

type localOptions struct {
	size int
}

type LocalCache[V any] struct {
	options  *localOptions
	lruCache *lru.Cache[string, localValue[V]]
}

type localValue[V any] struct {
	value      V
	expiration time.Time
}

type LocalOption func(*localOptions)

// WithSize sets the size of the cache default is 1000
func WithSize(size int) LocalOption {
	return func(options *localOptions) {
		options.size = size
	}
}

func NewLocalCache[V any](opts ...LocalOption) *LocalCache[V] {
	options := &localOptions{
		size: 1000,
	}
	for _, opt := range opts {
		opt(options)
	}

	lruCache, err := lru.New[string, localValue[V]](options.size)

	if err != nil {
		log.Panicf("local cache new lru err: %v", err)
	}

	return &LocalCache[V]{
		options:  options,
		lruCache: lruCache,
	}
}

func (l *LocalCache[V]) Get(ctx context.Context, key string) (V, error) {
	v, ok := l.lruCache.Get(key)
	if !ok {
		return v.value, ErrNoData
	}

	if v.expiration.Before(time.Now()) {
		// data is expired so remove it
		l.lruCache.Remove(key)
		return v.value, ErrExpire
	}
	return v.value, nil
}

func (l *LocalCache[V]) Set(ctx context.Context, key string, value V, expiration time.Duration) error {
	val := localValue[V]{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	l.lruCache.Add(key, val)
	return nil
}

func (l *LocalCache[V]) Delete(ctx context.Context, key string) error {
	l.lruCache.Remove(key)
	return nil
}
