package tasks

var systemTasks TaskManager

func Register(name string, task Task, config *TaskConfig) error {
	return systemTasks.Register(name, task, config)
}

func Close() error {
	return systemTasks.Close()
}
