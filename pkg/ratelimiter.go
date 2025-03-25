package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter interface {
	Allow(key string) (bool, error)
}

type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisLimiter(addr string, limit int, window time.Duration) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisRateLimiter) Allow(key string) (bool, error) {
	ctx := context.Background()
	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.client.Expire(ctx, key, rl.window)
	}

	return count <= int64(rl.limit), nil
}
