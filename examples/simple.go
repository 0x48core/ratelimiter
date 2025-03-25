package examples

import (
	"fmt"
	"time"

	ratelimiter "github.com/augustus281/ratelimiter/pkg"
)

func main() {
	limiter := ratelimiter.NewRedisLimiter("localhost:6379", 5, time.Minute)

	for i := 0; i < 7; i++ {
		allowed, _ := limiter.Allow("user:123")
		if allowed {
			fmt.Println("✅ Request allowed")
		} else {
			fmt.Println("❌ Too many requests")
		}
	}
}
