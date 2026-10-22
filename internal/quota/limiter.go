package quota

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	counts   map[string]int
	currDate string
}

func NewLimiter() *Limiter {
	return &Limiter{
		counts:   make(map[string]int),
		currDate: time.Now().UTC().Format("2006-01-02"),
	}
}

func (l *Limiter) Allow(ctx context.Context, service string, dailyCap int) (bool, int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().UTC().Format("2006-01-02")
	if today != l.currDate {
		l.currDate = today
		l.counts = make(map[string]int)
	}

	key := fmt.Sprintf("%s:%s", service, today)
	count := l.counts[key]
	if count >= dailyCap {
		return false, count, nil
	}

	l.counts[key] = count + 1
	return true, count + 1, nil
}
