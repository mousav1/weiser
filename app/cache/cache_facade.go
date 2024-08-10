package cache

import "time"

// CacheFacade provides a simplified interface to interact with the cache.
type CacheFacade struct{}

func (cf *CacheFacade) Get(key string, value interface{}) error {
	return GetCacheInstance().Get(key, value)
}

func (cf *CacheFacade) Set(key string, value interface{}, expiration time.Duration) error {
	return GetCacheInstance().Set(key, value, expiration)
}

func (cf *CacheFacade) Delete(key string) error {
	return GetCacheInstance().Delete(key)
}

func (cf *CacheFacade) Exists(key string) (bool, error) {
	return GetCacheInstance().Exists(key)
}

func (cf *CacheFacade) Flush() error {
	return GetCacheInstance().Flush()
}

func (cf *CacheFacade) Stats() (CacheStats, error) {
	return GetCacheInstance().Stats()
}
