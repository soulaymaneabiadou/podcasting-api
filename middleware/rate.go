package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	clientRequestCounts = make(map[string]map[string]int) // Map to store client request counts
	lastResetTime       time.Time                         // Last time the request counts were reset
	mutex               sync.Mutex
)

func RateLimit(maxReqs, perMin int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestPath := c.Request.URL.Path

		// Check if it's time to reset request counts
		if time.Since(lastResetTime).Minutes() >= float64(perMin) {
			mutex.Lock()
			clientRequestCounts = make(map[string]map[string]int)
			lastResetTime = time.Now()
			mutex.Unlock()
		}

		// Initialize the request count for the client and path if it doesn't exist
		mutex.Lock()
		if clientRequestCounts[clientIP] == nil {
			clientRequestCounts[clientIP] = make(map[string]int)
		}
		// Increment the request count for the client and path
		clientRequestCounts[clientIP][requestPath]++
		count := clientRequestCounts[clientIP][requestPath]
		mutex.Unlock()

		// Check if the request count exceeds the rate limit
		if count > maxReqs {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Continue processing the request
		c.Next()
	}
}
