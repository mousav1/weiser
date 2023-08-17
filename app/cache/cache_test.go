package cache

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache_Get(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	err = c.Get("nonexistent", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestInMemoryCache_Set(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestInMemoryCache_Delete(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	err = c.Delete("key")
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestInMemoryCache_Exists(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	exists, err := c.Exists("key")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = c.Exists("nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestInMemoryCache_Flush(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	err = c.Flush()
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestInMemoryCache_Stats(t *testing.T) {
	c := NewInMemoryCache(nil)
	value := "test value"
	err := c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	stats, err := c.Stats()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), stats.ItemsCount)
}

func TestFileCache_Get(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	err = c.Get("nonexistent", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestFileCache_Set(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestFileCache_Delete(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	err = c.Delete("key")
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestFileCache_Exists(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	exists, err := c.Exists("key")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = c.Exists("nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestFileCache_Flush(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	err = c.Flush()
	assert.NoError(t, err)

	var result string
	err = c.Get("key", &result)
	assert.Error(t, ErrCacheMiss, err)
}

func TestFileCache_Stats(t *testing.T) {
	// Create a temporary file for testing
	file, err := ioutil.TempFile("", "cache_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	c := NewFileCache(file.Name())
	value := "test value"
	err = c.Set("key", value, time.Minute)
	assert.NoError(t, err)

	stats, err := c.Stats()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), stats.ItemsCount)
}
