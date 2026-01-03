package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	windowStart time.Time
	windowLimit time.Duration

	tickets      int
	ticketsLimit int

	mu sync.Mutex
}

func NewLimiter(windowLimit time.Duration, ticketsLimit int) *Limiter {
	return &Limiter{
		windowStart:  time.Now(),
		windowLimit:  windowLimit,
		tickets:      0,
		ticketsLimit: ticketsLimit,
		mu:           sync.Mutex{},
	}
}

func (l *Limiter) Allow() bool {
	now := time.Now()

	if now.Sub(l.windowStart) >= l.windowLimit {
		l.windowStart = now
		l.tickets = 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.tickets == l.ticketsLimit {
		return false
	}

	l.tickets++

	return true
}
