# go-cache

When you need to cache data in memory, there are some options: 1. Use a global variable 2. Use a struct that depends on 
a specific cache instance or something like that. However, when you need to migrate to another cache, for example from
Redis to Memcached, you will need to change all the code that uses the cache. This package aims to solve this problem with
dependency inversion principle: DO NOT depend on concretions, depend on abstractions. There are 2 adapters currently:
localcache(implemented by LRU algorithm: https://github.com/hashicorp/golang-lru) and Redis. You can create your own adapter 
by implementing the Cacher interface.

## Example

```go
package main

import (
	"context"
	"fmt"

	cache "github.com/kimistar/go-cache"
	"github.com/kimistar/go-cache/adapter"
)

type User struct {
	Cacher cache.Cacher
}

// NewUser creates a new User instance provider that depends on a cacher instance 
// provided somewhere else(maybe by wire) by the adapter package
func NewUser(c cache.Cacher) *User {
	return &User{Cacher: c}
}

func (u *User) GetName(ctx context.Context) {
	name, err := cache.Cache(ctx, u.Cacher, "yourkey", func() (string, error) {
		// Your business logic
		return "yourname", nil
	})
	if err != nil {
		// Handle error
	}
	// Do something with name
	fmt.Println(name)
}

func main() {
	// Create a new localcache instance
	c := adapter.NewLocalDefault()
	// Create a new User instance
	u := NewUser(c)
	u.GetName(context.Background())
}
```