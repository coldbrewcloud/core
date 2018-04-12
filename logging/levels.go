package logging

import _logging "github.com/op/go-logging"

type Level int

var (
	DEBUG    Level = Level(_logging.DEBUG)
	INFO           = Level(_logging.INFO)
	WARNING        = Level(_logging.WARNING)
	ERROR          = Level(_logging.ERROR)
)
