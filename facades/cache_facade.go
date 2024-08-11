package facades

import (
	"time"

	"github.com/mousav1/weiser/app/cache"
)

// CacheFacade provides a simplified interface to interact with the cache.
type CacheFacade struct{}

func (cf *CacheFacade) Get(key string, value interface{}) error {
	return cache.GetCacheInstance().Get(key, value)
}

func (cf *CacheFacade) Set(key string, value interface{}, expiration time.Duration) error {
	return cache.GetCacheInstance().Set(key, value, expiration)
}

func (cf *CacheFacade) Delete(key string) error {
	return cache.GetCacheInstance().Delete(key)
}

func (cf *CacheFacade) Exists(key string) (bool, error) {
	return cache.GetCacheInstance().Exists(key)
}

func (cf *CacheFacade) Flush() error {
	return cache.GetCacheInstance().Flush()
}

func (cf *CacheFacade) Stats() (cache.CacheStats, error) {
	return cache.GetCacheInstance().Stats()
}

func NewCacheFacade() *CacheFacade {
	return &CacheFacade{}
}
