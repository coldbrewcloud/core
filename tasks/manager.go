package tasks

import (
	"sync"
	"time"
)

type Manager struct {
	ch chan struct{}
	wg sync.WaitGroup
}

func NewManager() *Manager {
	return &Manager{
		ch: make(chan struct{}),
	}
}

func (m *Manager) Start(task Task) {
	if task == nil {
		return
	}

	go func() {
		m.wg.Add(1)
		defer m.wg.Done()

		var ok bool
		nextDelay := time.Duration(0)

		for {
			select {
			case <-m.ch:
				return
			case <-time.After(nextDelay):
			}

			nextDelay, ok = task()
			if !ok {
				return
			}
		}
	}()
}

func (m *Manager) StopAll() {
	// signal all tasks to stop
	close(m.ch)
	m.wg.Wait()
}
