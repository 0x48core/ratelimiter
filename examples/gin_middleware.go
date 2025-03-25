package examples

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	ratelimiter "github.com/augustus281/ratelimiter/pkg"
)

var limiter = ratelimiter.NewRedisLimiter("localhost:6379", 10, time.Minute)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		allowed, _ := limiter.Allow(key)
		if !allowed {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func mainExample() {
	r := gin.Default()
	r.Use(RateLimitMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	r.Run(":8080")
}
