package tasks

import "time"

type Task func() (time.Duration, bool)
