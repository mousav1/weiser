package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache is a global instance of the cache
var Cache *cache.Cache

func init() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

// Get retrieves an item from the cache
func Get(key string) (interface{}, bool) {
	return Cache.Get(key)
}

// Set sets an item in the cache
func Set(key string, value interface{}, expiration time.Duration) {
	Cache.Set(key, value, expiration)
}

// Delete removes an item from the cache
func Delete(key string) {
	Cache.Delete(key)
}
