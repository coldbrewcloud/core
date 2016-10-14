package tasks_test

import (
	"testing"
	"time"

	"github.com/coldbrewcloud/coldbrew-core/tasks"
	"github.com/stretchr/testify/assert"
)

func TestSystemTaskManager(t *testing.T) {
	// normal case
	ch := make(chan bool)
	err := tasks.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Second,
		RunParam:       ch,
	})
	assert.Nil(t, err)
	err = wait(ch, 1, 5*time.Second)
	assert.Nil(t, err)
	err = tasks.Close()
	assert.Nil(t, err)

	// one blocking task should not block other tasks
	ch1 := make(chan bool)
	err = tasks.Register("t1", ConsumerTask, &tasks.TaskConfig{
		IterationDelay:  time.Millisecond,
		ContinueOnError: true,
		RunParam:        ch1,
	}) // this consumer task will block until we send something to "ch1"
	assert.Nil(t, err)
	ch2 := make(chan bool)
	err = tasks.Register("t2", ProducerTask, &tasks.TaskConfig{
		IterationDelay:  time.Millisecond,
		ContinueOnError: true,
		RunParam:        ch2,
	})
	assert.Nil(t, err)
	err = wait(ch2, 100, 5*time.Second)
	assert.Nil(t, err)
	ch1 <- true // unblock task "t1"
	err = tasks.Close()
	assert.Nil(t, err)

	// registering the same task name
	ch1 = make(chan bool)
	err = tasks.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Millisecond,
		RunParam:       ch1,
	})
	assert.Nil(t, err)
	ch2 = make(chan bool)
	err = tasks.Register("t1", ProducerTask, &tasks.TaskConfig{
		IterationDelay: time.Millisecond,
		RunParam:       ch2,
	})
	assert.NotNil(t, err) // should fail: duplicate task name
	err = wait(ch1, 1, time.Second)
	assert.Nil(t, err)
	err = tasks.Close()
	assert.Nil(t, err)

	// Register() param validations
	ch = make(chan bool)
	err = tasks.Register("", ProducerTask, nil)
	assert.NotNil(t, err)
	err = tasks.Register("t1", nil, nil)
	assert.NotNil(t, err)
	err = tasks.Close()
	assert.Nil(t, err)
}
