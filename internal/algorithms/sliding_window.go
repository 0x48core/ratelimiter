package algorithms

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type SlidingWindow struct {
	client *redis.Client
	window time.Duration
	limit  int
	mu     sync.Mutex
}

func NewSlidingWindow(client *redis.Client, window time.Duration, limit int) *SlidingWindow {
	return &SlidingWindow{
		client: client,
		window: window,
		limit:  limit,
	}
}

func (sw *SlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now().UnixMilli()
	redisKey := fmt.Sprintf("sw:%s", key)

	windowStart := now - sw.window.Milliseconds()

	pipe := sw.client.TxPipeline()

	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart))

	countCmd := pipe.ZCard(ctx, redisKey)

	pipe.ZAdd(ctx, redisKey, redis.Z{Score: float64(now), Member: now})

	pipe.Expire(ctx, redisKey, sw.window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if countCmd.Val() >= int64(sw.limit) {
		return false, nil
	}

	return true, nil
}
