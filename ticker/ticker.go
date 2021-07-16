package ticker

import (
	"sync"
	"time"
)

type Ticker struct {
	*time.Ticker
	Done    chan struct{}
	stopped bool
	mutex   sync.Mutex
}

func (t *Ticker) Stop() {
	t.mutex.Lock()
	if !t.stopped {
		t.Ticker.Stop()
		close(t.Done)
		t.stopped = true
	}
	t.mutex.Unlock()
}

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{
		Ticker: time.NewTicker(d),
		Done:   make(chan struct{}),
	}
}
