package tasks

import "time"

type TaskConfig struct {
	IterationDelay  time.Duration
	ContinueOnError bool
	RunParam        interface{}
}

func DefaultTaskConfig() *TaskConfig {
	return &TaskConfig{
		IterationDelay:  time.Second,
		ContinueOnError: true,
		RunParam:        nil,
	}
}
