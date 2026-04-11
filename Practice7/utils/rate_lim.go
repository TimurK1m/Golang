package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	count     int
	timestamp time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
	limit   = 5
	window  = time.Minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// 👇 если есть JWT
		if userID, exists := c.Get("userID"); exists {
			key = userID.(string)
		} else {
			key = c.ClientIP()
		}

		mu.Lock()
		defer mu.Unlock()

		cl, exists := clients[key]

		if !exists || time.Since(cl.timestamp) > window {
			clients[key] = &client{count: 1, timestamp: time.Now()}
		} else {
			cl.count++
			if cl.count > limit {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "too many requests",
				})
				return
			}
		}

		c.Next()
	}
}