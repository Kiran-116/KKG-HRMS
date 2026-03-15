package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	window   time.Duration // time window
}

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	authLimiter = &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     5,           // 5 requests per minute
		window:   1 * time.Minute,
	}
	apiLimiter = &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     100,          // 100 requests per minute
		window:   1 * time.Minute,
	}
)

// RateLimitAuth limits authentication endpoints to 5 requests per minute
func RateLimitAuth() gin.HandlerFunc {
	return rateLimitMiddleware(authLimiter)
}

// RateLimitAPI limits general API endpoints to 100 requests per minute
func RateLimitAPI() gin.HandlerFunc {
	return rateLimitMiddleware(apiLimiter)
}

func rateLimitMiddleware(limiter *rateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		v, exists := limiter.visitors[ip]
		if !exists {
			limiter.visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			limiter.mu.Unlock()
			c.Next()
			return
		}

		// Reset if window has passed
		if time.Since(v.lastSeen) > limiter.window {
			v.count = 1
			v.lastSeen = time.Now()
			limiter.mu.Unlock()
			c.Next()
			return
		}

		// Check if limit exceeded
		if v.count >= limiter.rate {
			limiter.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		v.count++
		v.lastSeen = time.Now()
		limiter.mu.Unlock()

		c.Next()
	}
}

// CleanupVisitors periodically removes old visitors (runs in background)
func CleanupVisitors() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			cleanup(apiLimiter)
			cleanup(authLimiter)
		}
	}()
}

func cleanup(limiter *rateLimiter) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	for ip, v := range limiter.visitors {
		if time.Since(v.lastSeen) > 30*time.Minute {
			delete(limiter.visitors, ip)
		}
	}
}
