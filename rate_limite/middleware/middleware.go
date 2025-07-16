package middleware

import (
	"net/http"
	"strings"

	"github.com/dlcdev1/pos1/rate_limite/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(l limiter.LimiterInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := strings.TrimSpace(c.GetHeader("API_KEY"))

		allowed, err := l.Allow(token, ip)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to validate rate limit",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message":     "you have reached the maximum number of requests or actions allowed within a certain time frame",
				"ip":          ip,
				"status code": http.StatusTooManyRequests,
			})
			return
		}

		c.Next()
	}
}
