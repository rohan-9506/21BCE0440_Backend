package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// A struct to hold information about requests from a particular IP
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
}

// A struct to track the number of requests and when the first request happened
type Visitor struct {
	Requests    int
	LastRequest time.Time
}

var limiter = RateLimiter{
	visitors: make(map[string]*Visitor),
}

// RateLimitMiddleware is the middleware to apply rate limiting
func RateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	visitor, exists := limiter.visitors[ip]
	if !exists || time.Since(visitor.LastRequest) > time.Minute {
		// Reset or initialize a new visitor
		visitor = &Visitor{Requests: 1, LastRequest: time.Now()}
		limiter.visitors[ip] = visitor
	} else {
		visitor.Requests++
		if visitor.Requests > 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		visitor.LastRequest = time.Now()
	}

	c.Next()
}
