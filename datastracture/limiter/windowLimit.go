/*
 * @Author: zengzh
 * @Date: 2023-02-02 17:25:30
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-02-02 17:30:55
 */

// https://juejin.cn/post/7056068978862456846
package limiter

import (
	"errors"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit    int
	window   time.Duration
	counter  int
	lastTime time.Time
	mutex    sync.Mutex
}

type SlidingWindowLimiter struct {
	limit        int
	window       int64
	smallWindow  int64
	smallWindows int64
	counters     map[int64]int
	mutex        sync.Mutex
}

func NewFixWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:    limit,
		window:   window,
		lastTime: time.Now(),
	}
}

func (l *FixedWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	now := time.Now()
	if now.Sub(l.lastTime) > l.window {
		l.counter = 0
		l.lastTime = now
	}
	if l.counter >= l.limit {
		return false
	}
	l.counter++
	return true
}

func NewSlidingWindowLimiter(limit int, window time.Duration, smallWindow time.Duration) (*SlidingWindowLimiter, error) {
	if window%smallWindow != 0 {
		return nil, errors.New("window must be a multiple of smallWindow")
	}
	return &SlidingWindowLimiter{
		limit:        limit,
		window:       int64(window),
		smallWindow:  int64(smallWindow),
		smallWindows: int64(window / smallWindow),
		counters:     make(map[int64]int),
	}, nil
}

func (l *SlidingWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	currentSmallWindow := time.Now().UnixNano() / l.smallWindow * l.smallWindow
	startSmallWindow := currentSmallWindow - l.smallWindow*(l.smallWindows-1)
	var count int
	for smallWindow, counter := range l.counters {
		if smallWindow < startSmallWindow {
			delete(l.counters, smallWindow)
		} else {
			count += counter
		}
	}

	if count >= l.limit {
		return false
	}
	l.counters[currentSmallWindow]++
	return true
}
