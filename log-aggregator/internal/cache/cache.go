package cache

import (
	"fmt"
	"sync"
	"time"
)

// Item represents a cached item with expiration
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache interface defines the methods that any cache implementation must provide
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// InMemoryCache implements the Cache interface using sync.Map
type InMemoryCache struct {
	store sync.Map
}

// NewInMemoryCache creates a new instance of InMemoryCache
func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{}
	// Start the cleanup routine
	go cache.startCleanup()
	return cache
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	value, exists := c.store.Load(key)
	if !exists {
		return nil, false
	}

	item := value.(Item)
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.Delete(key)
		return nil, false
	}

	return item.Value, true
}

// Set adds a value to the cache with an optional TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}
	c.store.Store(key, Item{
		Value:      value,
		Expiration: expiration,
	})
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(key string) {
	c.store.Delete(key)
}

// Clear removes all items from the cache
func (c *InMemoryCache) Clear() {
	c.store = sync.Map{}
}

// startCleanup starts a goroutine that periodically removes expired items
func (c *InMemoryCache) startCleanup() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now().UnixNano()
			c.store.Range(func(key, value interface{}) bool {
				item := value.(Item)
				if item.Expiration > 0 && now > item.Expiration {
					c.store.Delete(key)
				}
				return true
			})
		}
	}()
}

// CacheKeyGenerator generates cache keys for different types of queries
type CacheKeyGenerator struct{}

// NewCacheKeyGenerator creates a new instance of CacheKeyGenerator
func NewCacheKeyGenerator() *CacheKeyGenerator {
	return &CacheKeyGenerator{}
}

// ForLog generates a cache key for a single log
func (g *CacheKeyGenerator) ForLog(id string) string {
	return "log:" + id
}

// ForLogList generates a cache key for a paginated log list
func (g *CacheKeyGenerator) ForLogList(limit, offset int) string {
	return fmt.Sprintf("logs:list:%d:%d", limit, offset)
}

// ForTimeRange generates a cache key for time range queries
func (g *CacheKeyGenerator) ForTimeRange(start, end time.Time, limit, offset int) string {
	return fmt.Sprintf("logs:range:%d:%d:%s:%s",
		limit, offset, start.Format(time.RFC3339), end.Format(time.RFC3339))
}
