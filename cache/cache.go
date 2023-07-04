package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

// Cache is a global instance of the cache
var Cache *cache.Cache

func init() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

// Get retrieves an item from the cache
func Get(key string) (interface{}, error) {
	value, found := Cache.Get(key)
	if !found {
		return nil, errors.New("cache: key not found")
	}

	return value, nil
}

// Set sets an item in the cache
func Set(key string, value interface{}, expiration time.Duration) error {
	Cache.Set(key, value, expiration)
	return nil
}

// Delete removes an item from the cache
func Delete(key string) error {
	Cache.Delete(key)
	return nil
}
