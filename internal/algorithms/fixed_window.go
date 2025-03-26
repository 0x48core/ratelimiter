package algorithms

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedWindow struct {
	client *redis.Client
	window time.Duration
	limit  int
	mu     sync.Mutex
}

func NewFixedWindow(client *redis.Client, window time.Duration, limit int) *FixedWindow {
	return &FixedWindow{
		client: client,
		window: window,
		limit:  limit,
	}
}

func (fw *FixedWindow) Allow(ctx context.Context, key string) (bool, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	currentTime := time.Now().Unix()

	windowStart := currentTime - (currentTime % int64(fw.window.Seconds()))

	redisKey := fmt.Sprintf("fw:%s:%d", key, windowStart)

	pipe := fw.client.TxPipeline()

	count := pipe.Incr(ctx, redisKey)

	pipe.Expire(ctx, redisKey, fw.window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if count.Val() > int64(fw.limit) {
		return false, nil
	}

	return true, nil
}
