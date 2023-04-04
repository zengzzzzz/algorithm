package limiter

import (
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	peakLevel     int
	currentLevel  int
	currentVelocy int
	lastTime      time.Time
	mutex         sync.Mutex
}

func NewLeakyBucketLimit(peakLevel, currentVelocy int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		peakLevel:     peakLevel,
		currentVelocy: currentVelocy,
		lastTime:      time.Now(),
	}
}

func (l *LeakyBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	now := time.Now()
	interval := now.Sub(l.lastTime)
	if interval >= time.Second {
		l.currentLevel = maxInt(0, l.currentLevel-int(interval/time.Second)*l.currentVelocy)
		l.lastTime = now
	}
	if l.currentLevel >= l.peakLevel {
		return false
	}
	l.currentLevel++
	return true
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
