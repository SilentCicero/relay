package model

import (
	"sync"
	"time"
)

const MaxAge = time.Minute * 20

type comm struct {
	accessed time.Time
	wSecret  string
	l        sync.RWMutex
	*Waiters
}

func newComm() comm {
	return comm{
		accessed: time.Now(),
		Waiters:  NewWaiters(),
	}
}

func (c *comm) WAuth(wSecret string) bool {
	c.l.Lock()
	defer c.l.Unlock()

	if c.wSecret == "" {
		c.wSecret = wSecret
	}
	return c.wSecret == wSecret
}

func (c *comm) Expired() bool {
	c.l.RLock()
	defer c.l.RUnlock()
	return time.Since(c.accessed) > MaxAge && !c.Waiters.HasWaiters()
}

func (c *comm) Accessed() {
	c.l.Lock()
	defer c.l.Unlock()
	c.accessed = time.Now()
}
