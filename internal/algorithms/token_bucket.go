package algorithms

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	client   *redis.Client
	key      string
	capacity int
	rate     float64
	mu       sync.Mutex
}

func NewTokenBucket(client *redis.Client, key string, capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		client:   client,
		key:      key,
		capacity: capacity,
		rate:     rate,
	}
}

func (tb *TokenBucket) Allow(ctx context.Context) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now().Unix()

	tokensStr, err := tb.client.HGet(ctx, tb.key, "tokens").Result()
	lastRefillStr, errLast := tb.client.HGet(ctx, tb.key, "last_refill_time").Result()

	var tokens int
	var lastRefillTime int64

	if err == redis.Nil || errLast == redis.Nil {
		tokens = tb.capacity
		lastRefillTime = now
	} else if err == nil && errLast == nil {
		tokens, _ = strconv.Atoi(tokensStr)
		lastRefillTime, _ = strconv.ParseInt(lastRefillStr, 10, 64)
	} else {
		return false, err
	}

	elapsed := now - lastRefillTime
	newTokens := Min(tb.capacity, tokens+int(float64(elapsed)*tb.rate))

	if newTokens > 0 {
		newTokens--
		_, err := tb.client.HSet(ctx, tb.key, "tokens", newTokens, "last_refill_time").Result()
		if err != nil {
			return false, err
		}
		tb.client.Expire(ctx, tb.key, 60*time.Second)
		return true, nil
	}

	return false, nil
}
