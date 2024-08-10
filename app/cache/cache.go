package cache

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// CacheDriver defines the interface for different cache drivers.
type CacheDriver interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
	Stats() (CacheStats, error)
}

// Cache is a high-level abstraction over CacheDriver.
type Cache struct {
	driver            CacheDriver
	defaultExpiration time.Duration
	onEvicted         func(key string, value interface{})
}

var (
	cacheInstance *Cache
	once          sync.Once
	ErrCacheMiss  = errors.New("cache: key not found")
)

// NewCache initializes a new Cache instance.
func InitializeCache(defaultExpiration time.Duration, onEvicted func(key string, value interface{})) error {
	var err error
	once.Do(func() {
		var driver CacheDriver
		cacheType := viper.GetString("cache.type")

		switch cacheType {
		case "redis":
			client := redis.NewClient(&redis.Options{
				Addr:     viper.GetString("cache.redis.addr"),
				Password: viper.GetString("cache.redis.password"),
				DB:       viper.GetInt("cache.redis.db"),
			})
			driver = NewRedisCache(client, defaultExpiration, onEvicted)
		case "file":
			filePath := viper.GetString("cache.file.path")
			driver = NewFileCache(filePath)
		case "memory":
			driver = NewInMemoryCache(onEvicted)
		default:
			err = errors.New("unsupported cache type")
			return
		}

		cacheInstance = &Cache{
			driver:            driver,
			defaultExpiration: defaultExpiration,
			onEvicted:         onEvicted,
		}
	})
	return err
}

// GetCache returns the initialized Cache instance.
func GetCacheInstance() *Cache {
	return cacheInstance
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string, value interface{}) error {
	return c.driver.Get(key, value)
}

// Set stores a value in the cache.
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}
	return c.driver.Set(key, value, expiration)
}

// Delete removes a value from the cache.
func (c *Cache) Delete(key string) error {
	return c.driver.Delete(key)
}

// Exists checks if a key exists in the cache.
func (c *Cache) Exists(key string) (bool, error) {
	return c.driver.Exists(key)
}

// Flush clears all items from the cache.
func (c *Cache) Flush() error {
	return c.driver.Flush()
}

// Stats returns cache statistics.
func (c *Cache) Stats() (CacheStats, error) {
	return c.driver.Stats()
}

// InMemoryCache is an in-memory cache implementation.
type InMemoryCache struct {
	cache     map[string][]byte
	expiry    map[string]time.Time
	mutex     sync.RWMutex
	onEvicted func(string, interface{})
}

// NewInMemoryCache initializes a new in-memory cache.
func NewInMemoryCache(onEvicted func(key string, value interface{})) *InMemoryCache {
	return &InMemoryCache{
		cache:     make(map[string][]byte),
		expiry:    make(map[string]time.Time),
		mutex:     sync.RWMutex{},
		onEvicted: onEvicted,
	}
}

// Get retrieves a value from the in-memory cache.
func (c *InMemoryCache) Get(key string, value interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if data, ok := c.cache[key]; ok {
		return json.Unmarshal(data, value)
	}

	return ErrCacheMiss
}

// Set stores a value in the in-memory cache.
func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = time.Minute
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = data
	c.expiry[key] = time.Now().Add(expiration)

	go c.startJanitor()

	return nil
}

// Delete removes a value from the in-memory cache.
func (c *InMemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.cache[key]; ok {
		delete(c.cache, key)
		delete(c.expiry, key)
		return nil
	}

	return ErrCacheMiss
}

// Exists checks if a key exists in the in-memory cache.
func (c *InMemoryCache) Exists(key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, ok := c.cache[key]
	return ok, nil
}

// Flush clears all items from the in-memory cache.
func (c *InMemoryCache) Flush() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string][]byte)
	c.expiry = make(map[string]time.Time)

	return nil
}

// Stats returns statistics of the in-memory cache.
func (c *InMemoryCache) Stats() (CacheStats, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return CacheStats{ItemsCount: int64(len(c.cache))}, nil
}

// startJanitor periodically cleans up expired items.
func (c *InMemoryCache) startJanitor() {
	for {
		time.Sleep(time.Minute)

		c.mutex.Lock()
		for key, expiry := range c.expiry {
			if time.Now().After(expiry) {
				if c.onEvicted != nil {
					data := c.cache[key]
					c.onEvicted(key, data)
					delete(c.cache, key)
					delete(c.expiry, key)
				}
			}
		}
		c.mutex.Unlock()
	}
}

// RedisCache is a Redis cache implementation.
type RedisCache struct {
	client            *redis.Client
	defaultExpiration time.Duration
	onEvicted         func(key string, value interface{})
}

// NewRedisCache initializes a new Redis cache.
func NewRedisCache(client *redis.Client, defaultExpiration time.Duration, onEvicted func(key string, value interface{})) *RedisCache {
	return &RedisCache{
		client:            client,
		defaultExpiration: defaultExpiration,
		onEvicted:         onEvicted,
	}
}

// Get retrieves a value from the Redis cache.
func (r *RedisCache) Get(key string, value interface{}) error {
	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(data, value)
}

// Set stores a value in the Redis cache.
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = r.defaultExpiration
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, data, expiration).Err()
}

// Delete removes a value from the Redis cache.
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// Exists checks if a key exists in the Redis cache.
func (r *RedisCache) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// Flush clears all items from the Redis cache.
func (r *RedisCache) Flush() error {
	return r.client.FlushDB(context.Background()).Err()
}

// Stats returns statistics of the Redis cache.
func (r *RedisCache) Stats() (CacheStats, error) {
	stats := CacheStats{}
	keys, err := r.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return stats, err
	}
	stats.ItemsCount = int64(len(keys))
	return stats, nil
}

// CacheItem represents an item stored in the file cache.
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// FileCache is a file-based cache implementation.
type FileCache struct {
	cache map[string]CacheItem
	mutex sync.Mutex
	path  string
}

// NewFileCache initializes a new file-based cache.
func NewFileCache(path string) *FileCache {
	return &FileCache{
		cache: make(map[string]CacheItem),
		path:  path,
	}
}

// Get retrieves a value from the file cache.
func (f *FileCache) Get(key string, value interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Read the entire cache file
	bytes, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}

	// Unmarshal the file content into the cache map
	if err := json.Unmarshal(bytes, &f.cache); err != nil {
		return err
	}

	// Get the cache item
	item, ok := f.cache[key]
	if !ok || item.ExpiresAt.Before(time.Now()) {
		return ErrCacheMiss
	}

	// Unmarshal the item.Value which is stored as a []byte
	data, ok := item.Value.([]byte)
	if !ok {
		return errors.New("cache: invalid data type")
	}

	return json.Unmarshal(data, value)
}

// Set stores a value in the file cache.
func (f *FileCache) Set(key string, value interface{}, expiration time.Duration) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.cache[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(expiration),
	}

	bytes, err := json.Marshal(f.cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.path, bytes, 0644)
}

// Delete removes a value from the file cache.
func (f *FileCache) Delete(key string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	delete(f.cache, key)

	bytes, err := json.Marshal(f.cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.path, bytes, 0644)
}

// Exists checks if a key exists in the file cache.
func (f *FileCache) Exists(key string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, ok := f.cache[key]
	return ok, nil
}

// Flush clears all items from the file cache.
func (f *FileCache) Flush() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.cache = make(map[string]CacheItem)

	bytes, err := json.Marshal(f.cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.path, bytes, 0644)
}

// Stats returns statistics of the file cache.
func (f *FileCache) Stats() (CacheStats, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return CacheStats{ItemsCount: int64(len(f.cache))}, nil
}

// CacheStats holds statistics about the cache.
type CacheStats struct {
	ItemsCount int64
}
