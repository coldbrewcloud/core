package debug

import (
	"os"
	"strings"
)

const (
	EnvKey = "DEBUG"
)

func Enabled() bool {
	switch strings.ToLower(os.Getenv(EnvKey)) {
	case "1", "true":
		return true
	}

	return false
}
