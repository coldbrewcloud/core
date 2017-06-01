package tasks

import "sync"

var DefaultManager *Manager
var initDefaultManager sync.Once

func Start(task Task) {
	initDefaultManager.Do(func() {
		DefaultManager = NewManager()
	})

	DefaultManager.Start(task)
}

func StopAll() {
	DefaultManager.StopAll()
}
