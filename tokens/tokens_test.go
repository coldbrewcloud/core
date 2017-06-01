package tokens_test

import (
	"testing"

	"github.com/coldbrewcloud/core/tokens"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	token := tokens.New()
	assert.NotEmpty(t, token)
}
