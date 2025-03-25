package algorithms

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity  int
	leakyRate time.Duration
	tokens    int
	lastTime  time.Time
	mu        sync.Mutex
}

func NewLeakyBucket(capacity int, leakyRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:  capacity,
		leakyRate: leakyRate,
		tokens:    capacity,
		lastTime:  time.Now(),
	}
}

func (lb *LeakyBucket) AllowRequest() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	leakedTime := now.Sub(lb.lastTime)

	newTokens := int(leakedTime / lb.leakyRate)
	if newTokens > 0 {
		lb.tokens = Min(lb.capacity, lb.tokens+newTokens)
		lb.lastTime = now
	}

	if lb.tokens > 0 {
		lb.tokens--
		return true
	}

	return false
}
