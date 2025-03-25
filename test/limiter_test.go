package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ratelimiter "github.com/augustus281/ratelimiter/pkg"
)

func TestRateLimiter(t *testing.T) {
	limiter := ratelimiter.NewRedisLimiter("localhost:6379", 2, time.Second)

	allowed, _ := limiter.Allow("test_user")
	assert.True(t, allowed)

	allowed, _ = limiter.Allow("test_user")
	assert.True(t, allowed)

	allowed, _ = limiter.Allow("test_user")
	assert.False(t, allowed)
}
