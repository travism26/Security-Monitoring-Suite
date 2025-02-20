package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheKeyGenerator_UserIsolation(t *testing.T) {
	generator := NewCacheKeyGenerator()

	t.Run("ForLog generates unique keys per user", func(t *testing.T) {
		user1Key := generator.ForLog("user1", "log123")
		user2Key := generator.ForLog("user2", "log123")

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})

	t.Run("ForLogList generates unique keys per user", func(t *testing.T) {
		user1Key := generator.ForLogList("user1", 10, 0)
		user2Key := generator.ForLogList("user2", 10, 0)

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})

	t.Run("ForTimeRange generates unique keys per user", func(t *testing.T) {
		now := time.Now()
		start := now.Add(-1 * time.Hour)
		end := now

		user1Key := generator.ForTimeRange("user1", start, end, 10, 0)
		user2Key := generator.ForTimeRange("user2", start, end, 10, 0)

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})

	t.Run("ForTimeRangeCount generates unique keys per user", func(t *testing.T) {
		now := time.Now()
		start := now.Add(-1 * time.Hour)
		end := now

		user1Key := generator.ForTimeRangeCount("user1", start, end)
		user2Key := generator.ForTimeRangeCount("user2", start, end)

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})

	t.Run("ForHostLogs generates unique keys per user", func(t *testing.T) {
		user1Key := generator.ForHostLogs("user1", "host1", 10, 0)
		user2Key := generator.ForHostLogs("user2", "host1", 10, 0)

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})

	t.Run("ForLevelLogs generates unique keys per user", func(t *testing.T) {
		user1Key := generator.ForLevelLogs("user1", "ERROR", 10, 0)
		user2Key := generator.ForLevelLogs("user2", "ERROR", 10, 0)

		assert.NotEqual(t, user1Key, user2Key, "Cache keys for different users should be unique")
		assert.Contains(t, user1Key, "user1", "Cache key should contain user1 ID")
		assert.Contains(t, user2Key, "user2", "Cache key should contain user2 ID")
	})
}

func TestInMemoryCache_UserIsolation(t *testing.T) {
	t.Run("Different users cannot access each other's cached data", func(t *testing.T) {
		cache := NewInMemoryCache()
		generator := NewCacheKeyGenerator()

		// Setup test data
		user1Data := map[string]string{"key": "user1-value"}
		user2Data := map[string]string{"key": "user2-value"}

		// Cache data for both users
		user1Key := generator.ForLog("user1", "log123")
		user2Key := generator.ForLog("user2", "log123")

		cache.Set(user1Key, user1Data, time.Minute)
		cache.Set(user2Key, user2Data, time.Minute)

		// Verify user1 can only access their data
		value1, exists1 := cache.Get(user1Key)
		assert.True(t, exists1, "User1's data should exist in cache")
		assert.Equal(t, user1Data, value1, "User1 should get their own data")

		// Verify user2 can only access their data
		value2, exists2 := cache.Get(user2Key)
		assert.True(t, exists2, "User2's data should exist in cache")
		assert.Equal(t, user2Data, value2, "User2 should get their own data")

		// Verify users cannot access each other's data
		wrongKey1 := generator.ForLog("user1", "log456")
		wrongKey2 := generator.ForLog("user2", "log456")

		_, exists3 := cache.Get(wrongKey1)
		assert.False(t, exists3, "User1 should not access non-existent data")

		_, exists4 := cache.Get(wrongKey2)
		assert.False(t, exists4, "User2 should not access non-existent data")
	})

	t.Run("Cache expiration works correctly per user", func(t *testing.T) {
		cache := NewInMemoryCache()
		generator := NewCacheKeyGenerator()
		user1Data := "user1-data"
		user2Data := "user2-data"

		user1Key := generator.ForLog("user1", "expiring-log")
		user2Key := generator.ForLog("user2", "expiring-log")

		// Set data with different expiration times
		cache.Set(user1Key, user1Data, 50*time.Millisecond)
		cache.Set(user2Key, user2Data, time.Minute)

		// Verify both exist initially
		_, exists1 := cache.Get(user1Key)
		_, exists2 := cache.Get(user2Key)
		assert.True(t, exists1, "User1's data should exist initially")
		assert.True(t, exists2, "User2's data should exist initially")

		// Wait for user1's data to expire
		time.Sleep(100 * time.Millisecond)

		// Verify user1's data expired but user2's remains
		_, exists3 := cache.Get(user1Key)
		_, exists4 := cache.Get(user2Key)
		assert.False(t, exists3, "User1's data should have expired")
		assert.True(t, exists4, "User2's data should still exist")
	})

	t.Run("Cache clear affects all users", func(t *testing.T) {
		cache := NewInMemoryCache()
		generator := NewCacheKeyGenerator()
		// Set data for multiple users
		cache.Set(generator.ForLog("user1", "log1"), "data1", time.Minute)
		cache.Set(generator.ForLog("user2", "log2"), "data2", time.Minute)

		// Clear cache
		cache.Clear()

		// Verify no user can access previous data
		_, exists1 := cache.Get(generator.ForLog("user1", "log1"))
		_, exists2 := cache.Get(generator.ForLog("user2", "log2"))

		assert.False(t, exists1, "User1's data should be cleared")
		assert.False(t, exists2, "User2's data should be cleared")
	})

	t.Run("Cache delete affects only specific user data", func(t *testing.T) {
		cache := NewInMemoryCache()
		generator := NewCacheKeyGenerator()
		user1Key := generator.ForLog("user1", "log1")
		user2Key := generator.ForLog("user2", "log2")

		// Set data for both users
		cache.Set(user1Key, "data1", time.Minute)
		cache.Set(user2Key, "data2", time.Minute)

		// Delete only user1's data
		cache.Delete(user1Key)

		// Verify user1's data is gone but user2's remains
		_, exists1 := cache.Get(user1Key)
		_, exists2 := cache.Get(user2Key)

		assert.False(t, exists1, "User1's data should be deleted")
		assert.True(t, exists2, "User2's data should still exist")
	})
}
