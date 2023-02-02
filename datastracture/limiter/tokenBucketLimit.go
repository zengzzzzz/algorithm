package limiter

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	capacity      int
	currentTokens int
	rate          int
	lastTime      time.Time
	mutex         sync.Mutex
}

func NewTokenBucketLimiter(capacity, rate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (l *TokenBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	now := time.Now()
	interval := now.Sub(l.lastTime)
	if interval >= time.Second {
		l.currentTokens = minInt(l.capacity, l.currentTokens+int(interval/time.Second)*l.rate)
		l.lastTime = now
	}
	if l.currentTokens == 0 {
		return false
	}
	l.currentTokens--
	return true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
