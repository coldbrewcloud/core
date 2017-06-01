package tasks_test

import (
	"testing"
	"time"

	"github.com/coldbrewcloud/core/tasks"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	m := tasks.NewManager()
	assert.NotNil(t, m)
}

func TestManager_Start(t *testing.T) {
	// normal case
	m := tasks.NewManager()
	ch := make(chan struct{})
	m.Start(producer(ch))
	assertWait(t, ch, 1, 5*time.Second)
	m.StopAll()

	// one blocking task should not block other tasks
	m = tasks.NewManager()
	ch1 := make(chan struct{})
	m.Start(consumer(ch1)) // this consumer task will block until we send something to "ch1"
	ch2 := make(chan struct{})
	m.Start(producer(ch2))
	assertWait(t, ch2, 100, 5*time.Second)
	ch1 <- struct{}{} // unblock task "t1"
	m.StopAll()

	// task == nil : nothing will run
	m = tasks.NewManager()
	ch = make(chan struct{})
	m.Start(nil)
	m.StopAll()
}

func TestManager_StopAll(t *testing.T) {
	m := tasks.NewManager()

	// normal case
	ch := make(chan struct{})
	m.Start(producer(ch))
	assertWait(t, ch, 1, 5*time.Second)
	m.StopAll()

	// make sure producer is not running any more
	select {
	case <-ch:
		assert.Fail(t, "producer still running")
	default:
	}
}

func producer(ch chan struct{}) tasks.Task {
	return func() (time.Duration, bool) {
		ch <- struct{}{}
		return 1 * time.Millisecond, true
	}
}

func consumer(ch chan struct{}) tasks.Task {
	return func() (time.Duration, bool) {
		select {
		case _, ok := <-ch:
			if ok {
				return 1 * time.Millisecond, true
			}

			return 0, false
		}
	}
}

// wait until channel "ch" receives "expected" signals or times out at "timeout"
func assertWait(t *testing.T, ch chan struct{}, expected int, timeout time.Duration) {
	cnt := 0
	deadline := time.Now().Add(timeout)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				assert.Fail(t, "channel closed")
			}

			cnt++
			if cnt == expected {
				return
			}
		default:
			if time.Now().After(deadline) {
				assert.Fail(t, "timed out")
			}
		}
	}
}
