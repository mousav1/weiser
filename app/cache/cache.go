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

type CacheDriver interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
	Stats() (CacheStats, error)
}

type Cache struct {
	driver            CacheDriver
	defaultExpiration time.Duration
	onEvicted         func(key string, value interface{})
}

var ErrCacheMiss = errors.New("cache: key not found")
var cache *Cache

func NewCache(defaultExpiration time.Duration, onEvicted func(key string, value interface{})) error {
	var driver CacheDriver
	cacheType := viper.GetString("cache.type")

	switch cacheType {
	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		driver = NewRedisCache(client, defaultExpiration, onEvicted)
	case "file":
		filePath := viper.GetString("cache.file.path")
		driver = NewFileCache(filePath)
	case "memory":
		driver = NewInMemoryCache(onEvicted)
	default:
		panic("unsupported cache type")
	}

	cache = &Cache{
		driver:            driver,
		defaultExpiration: defaultExpiration,
		onEvicted:         onEvicted,
	}
	return nil
}

// GetSessionManager returns the initialized session manager
func GetCache() *Cache {
	return cache
}

func (c *Cache) Get(key string, value interface{}) error {
	return c.driver.Get(key, value)
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}
	return c.driver.Set(key, value, expiration)
}

func (c *Cache) Delete(key string) error {
	return c.driver.Delete(key)
}

func (c *Cache) Exists(key string) (bool, error) {
	return c.driver.Exists(key)
}

func (c *Cache) Flush() error {
	return c.driver.Flush()
}

func (c *Cache) Stats() (CacheStats, error) {
	stats, err := c.driver.Stats()
	if err != nil {
		return CacheStats{}, err
	}
	return stats, nil
}

type inMemoryCache struct {
	cache     map[string][]byte
	expiry    map[string]time.Time
	mutex     sync.RWMutex
	onEvicted func(string, interface{})
}

func NewInMemoryCache(onEvicted func(key string, value interface{})) *inMemoryCache {
	return &inMemoryCache{
		cache:     make(map[string][]byte),
		expiry:    make(map[string]time.Time),
		mutex:     sync.RWMutex{},
		onEvicted: onEvicted,
	}
}

func (c *inMemoryCache) Get(key string, value interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if data, ok := c.cache[key]; ok {
		return json.Unmarshal(data, value)
	}

	return ErrCacheMiss
}

func (c *inMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = time.Minute
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		c.cache = make(map[string][]byte)
		c.expiry = make(map[string]time.Time)
	}

	c.cache[key] = data
	c.expiry[key] = time.Now().Add(expiration)

	go c.startJanitor()

	return nil
}

func (c *inMemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.cache[key]; ok {
		delete(c.cache, key)
		delete(c.expiry, key)
		return nil
	}

	return ErrCacheMiss
}

func (c *inMemoryCache) Exists(key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.cache[key]; ok {
		return true, nil
	}

	return false, nil
}

func (c *inMemoryCache) Flush() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string][]byte)
	c.expiry = make(map[string]time.Time)

	return nil
}

func (c *inMemoryCache) Stats() (CacheStats, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stats := CacheStats{
		ItemsCount: int64(len(c.cache)),
	}
	return stats, nil
}

func (c *inMemoryCache) startJanitor() {
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

type RedisCache struct {
	client            *redis.Client
	defaultExpiration time.Duration
	onEvicted         func(key string, value interface{})
}

func NewRedisCache(client *redis.Client, defaultExpiration time.Duration, onEvicted func(key string, value interface{})) *RedisCache {
	return &RedisCache{
		client:            client,
		defaultExpiration: defaultExpiration,
		onEvicted:         onEvicted,
	}
}

func (c *RedisCache) Get(key string, value interface{}) error {
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(data, value)
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), key, data, expiration).Err()
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *RedisCache) Exists(key string) (bool, error) {
	exists, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (c *RedisCache) Flush() error {
	return c.client.FlushDB(context.Background()).Err()
}

func (c *RedisCache) Stats() (CacheStats, error) {
	stats := CacheStats{}
	keys, err := c.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return stats, err
	}
	stats.ItemsCount = int64(len(keys))
	return stats, nil
}

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type FileCache struct {
	cache map[string]CacheItem
	mutex sync.Mutex
	path  string
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		cache: make(map[string]CacheItem),
		path:  path,
	}
}

func (c *FileCache) Get(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	bytes, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &c.cache); err != nil {
		return err
	}

	item, ok := c.cache[key]
	if !ok || item.ExpiresAt.Before(time.Now()) {
		return ErrCacheMiss
	}

	jsonValue, err := json.Marshal(item.Value)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonValue, value); err != nil {
		return err
	}

	return nil
}

func (c *FileCache) Set(key string, value interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(expiration),
	}

	bytes, err := json.Marshal(c.cache)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (c *FileCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)

	bytes, err := json.Marshal(c.cache)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

type CacheStats struct {
	ItemsCount int64
}

func (c *FileCache) Exists(key string) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.cache[key]
	return ok, nil
}

func (c *FileCache) Flush() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]CacheItem)

	bytes, err := json.Marshal(c.cache)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (c *FileCache) Stats() (CacheStats, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stats := CacheStats{
		ItemsCount: int64(len(c.cache)),
	}
	return stats, nil
}
