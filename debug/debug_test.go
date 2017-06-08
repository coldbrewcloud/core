package debug_test

import (
	"os"
	"testing"

	"github.com/coldbrewcloud/core/debug"
	"github.com/stretchr/testify/assert"
)

func TestDebugEnabled(t *testing.T) {
	oldValue, restore := os.LookupEnv(debug.EnvKey)

	os.Unsetenv(debug.EnvKey)
	assert.False(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "1")
	assert.True(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "true")
	assert.True(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "TRUE")
	assert.True(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "True")
	assert.True(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "TrUe")
	assert.True(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "0")
	assert.False(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "false")
	assert.False(t, debug.Enabled())

	os.Setenv(debug.EnvKey, "truefalse")
	assert.False(t, debug.Enabled())

	if restore {
		os.Setenv(debug.EnvKey, oldValue)
	}
}
