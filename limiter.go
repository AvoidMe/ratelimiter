package ratelimiter

// This is a naive implementation of token bucket rate limiter
// https://en.wikipedia.org/wiki/Token_bucket

import (
	"context"
	"sync"
	"time"
)

type Limiter struct {
	ctx        context.Context
	cancel     context.CancelFunc
	tokens     float64
	limit      float64
	bucketSize float64
	interval   time.Duration
	lock       sync.Mutex
}

func NewLimiter(limit float64, interval time.Duration) *Limiter {
	ctx, cancel := context.WithCancel(context.Background())
	realLimit := limit
	if interval > time.Second {
		realLimit = limit / interval.Seconds()
		interval = time.Second
	}
	return &Limiter{
		ctx:        ctx,
		cancel:     cancel,
		tokens:     limit,
		limit:      realLimit,
		bucketSize: limit,
		interval:   interval,
	}
}

func (l *Limiter) Start() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.Refill()
			ticker.Reset(l.interval)
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *Limiter) Stop() {
	l.cancel()
}

func (l *Limiter) Refill() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.tokens = min(l.bucketSize, l.tokens+l.limit)
}

func (l *Limiter) AmountLeft() float64 {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.tokens
}

func (l *Limiter) Get() bool {
	return l.GetN(1)
}

func (l *Limiter) GetN(count int) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.tokens >= float64(count) {
		l.tokens -= float64(count)
		return true
	}
	return false
}
