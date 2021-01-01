package model

import (
	"sync"
)

type Waiters struct {
	waiters  int
	waitChan chan struct{}
	m        sync.RWMutex
}

func NewWaiters() *Waiters {
	wc := make(chan struct{})
	close(wc)
	return &Waiters{waitChan: wc}
}

func (w *Waiters) AddWaiter() {
	w.m.Lock()
	defer w.m.Unlock()

	if w.waiters == 0 {
		w.waitChan = make(chan struct{})
	}

	w.waiters++
}

func (w *Waiters) RemoveWaiter() {
	w.m.Lock()
	defer w.m.Unlock()

	if w.waiters > 0 {
		w.waiters--

		if w.waiters == 0 {
			close(w.waitChan)
		}
	}
}

func (w *Waiters) Wait() <-chan struct{} {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.waitChan
}

func (w *Waiters) WaiterCount() int {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.waiters
}

func (w *Waiters) HasWaiters() bool {
	return w.WaiterCount() > 0
}
