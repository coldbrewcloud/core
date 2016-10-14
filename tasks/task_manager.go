package tasks

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TaskManager struct {
	lock        sync.Mutex
	killSignals map[string]chan bool
	closing     bool
	waitGroup   sync.WaitGroup
}

func (tm *TaskManager) Register(name string, task Task, config *TaskConfig) error {
	if name == "" {
		return errors.New("name is empty")
	}
	if task == nil {
		return errors.New("task is nil")
	}
	if config == nil {
		config = DefaultTaskConfig()
	}

	tm.lock.Lock()
	defer tm.lock.Unlock()

	if tm.killSignals == nil {
		tm.killSignals = make(map[string]chan bool)
	}

	if _, exists := tm.killSignals[name]; exists {
		return fmt.Errorf("Worker name %q already registered", name)
	}

	killSignal := make(chan bool)
	tm.killSignals[name] = killSignal

	tm.waitGroup.Add(1)

	go func() {
		defer tm.waitGroup.Done()

		for !tm.closing {
			if err := task(config.RunParam); err != nil {
				if !config.ContinueOnError {
					break
				}
			}

			// cancellable waiting
		l1:
			for {
				select {
				case <-time.After(config.IterationDelay):
					break l1
				case <-killSignal:
					return // should stop
				}
			}
		}
	}()

	return nil
}

func (tm *TaskManager) Close() error {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	// signal all workers to stop
	tm.closing = true
	for _, s := range tm.killSignals {
		s <- true
	}

	// wait until all workers stop
	tm.waitGroup.Wait()

	// reset data too
	tm.killSignals = make(map[string]chan bool)
	tm.closing = false

	return nil
}
