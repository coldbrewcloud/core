package tasks_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/coldbrewcloud/coldbrew-core/tasks"
	"github.com/stretchr/testify/assert"
)

func ProducerTask(v interface{}) error {
	ch := v.(chan bool)
	ch <- true
	return nil
}

func ConsumerTask(v interface{}) error {
	ch := v.(chan bool)
	_, ok := <-ch
	if !ok {
		return fmt.Errorf("channel closed")
	}
	return nil
}

func TestTaskManager(t *testing.T) {
	manager := &tasks.TaskManager{}

	// normal case
	ch := make(chan bool)
	err := manager.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Second,
		RunParam:       ch,
	})
	assert.Nil(t, err)
	err = wait(ch, 1, 5*time.Second)
	assert.Nil(t, err)
	err = manager.Close()
	assert.Nil(t, err)

	// one blocking task should not block other tasks
	ch1 := make(chan bool)
	err = manager.Register("t1", ConsumerTask, &tasks.TaskConfig{
		IterationDelay:  time.Millisecond,
		ContinueOnError: true,
		RunParam:        ch1,
	}) // this consumer task will block until we send something to "ch1"
	assert.Nil(t, err)
	ch2 := make(chan bool)
	err = manager.Register("t2", ProducerTask, &tasks.TaskConfig{
		IterationDelay:  time.Millisecond,
		ContinueOnError: true,
		RunParam:        ch2,
	})
	assert.Nil(t, err)
	err = wait(ch2, 100, 5*time.Second)
	assert.Nil(t, err)
	ch1 <- true // unblock task "t1"
	err = manager.Close()
	assert.Nil(t, err)

	// registering the same task name
	ch1 = make(chan bool)
	err = manager.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Millisecond,
		RunParam:       ch1,
	})
	assert.Nil(t, err)
	ch2 = make(chan bool)
	err = manager.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Millisecond,
		RunParam:       ch2,
	})
	assert.NotNil(t, err) // should fail: duplicate task name
	err = wait(ch1, 1, time.Second)
	assert.Nil(t, err)
	err = manager.Close()
	assert.Nil(t, err)

	// Register() param validations
	ch = make(chan bool)
	err = manager.Register("", ProducerTask, nil)
	assert.NotNil(t, err)
	err = manager.Register("t1", nil, nil)
	assert.NotNil(t, err)
	err = manager.Close()
	assert.Nil(t, err)

	// Default config
	defaultConfig := tasks.DefaultTaskConfig()
	assert.NotNil(t, defaultConfig)
}

// wait until channel "ch" receives "numSignals" signals or times out at "timeout"
func wait(ch chan bool, numSignals int, timeout time.Duration) error {
	signals := 0
	startTime := time.Now()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return fmt.Errorf("Channel closed")
			}
			signals++
			if signals == numSignals {
				return nil
			}
		default:
			if time.Now().Sub(startTime) > timeout {
				return fmt.Errorf("Time out")
			}
		}
	}
}
