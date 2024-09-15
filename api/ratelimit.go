package api

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/gin-gonic/gin"
)

var rateLimiter = tollbooth.NewLimiter(100, &limiter.ExpirableOptions{
	DefaultExpirationTTL: time.Minute,
})

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpLimiter := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{
			DefaultExpirationTTL: time.Minute,
		})

		httpStatus := tollbooth.Limit(c.ClientIP(), httpLimiter, c.Writer, c.Request)
		if httpStatus != http.StatusOK {
			c.AbortWithStatus(httpStatus)
			return
		}

		c.Next()
	}
}
