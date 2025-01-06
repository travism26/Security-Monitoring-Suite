package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	sync.RWMutex
	requests    map[string][]time.Time
	windowSize  time.Duration
	maxRequests int
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(windowSize time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string][]time.Time),
		windowSize:  windowSize,
		maxRequests: maxRequests,
	}
}

// cleanup removes old requests outside the window
func (rl *RateLimiter) cleanup(key string, now time.Time) {
	windowStart := now.Add(-rl.windowSize)
	var validRequests []time.Time

	for _, timestamp := range rl.requests[key] {
		if timestamp.After(windowStart) {
			validRequests = append(validRequests, timestamp)
		}
	}

	if len(validRequests) == 0 {
		delete(rl.requests, key)
	} else {
		rl.requests[key] = validRequests
	}
}

// isAllowed checks if a request is allowed based on the rate limit
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	rl.cleanup(key, now)

	if len(rl.requests[key]) >= rl.maxRequests {
		return false
	}

	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger logs request details
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// Recovery recovers from panics
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// RateLimit creates a rate limiting middleware
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(window, requests)
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.isAllowed(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIKeyAuth creates an API key authentication middleware
func APIKeyAuth(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			c.Abort()
			return
		}

		isValid := false
		for _, key := range validAPIKeys {
			if strings.EqualFold(apiKey, key) {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
