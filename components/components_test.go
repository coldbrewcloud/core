package components_test

import (
	"testing"

	"github.com/coldbrewcloud/core/components"
	"github.com/stretchr/testify/assert"
)

func TestSystemComponents(t *testing.T) {
	// initial state
	c := components.Get(regularComponentName)
	assert.Nil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs := components.Close()
	assert.Nil(t, errs)

	// register regular component
	err := components.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	c = components.Get(regularComponentName)
	assert.NotNil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs = components.Close()
	assert.Nil(t, errs)

	// register init-error component
	err = components.Register(initErrorComponentName, &InitErrorComponent{})
	assert.NotNil(t, err)
	c = components.Get(regularComponentName)
	assert.Nil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs = components.Close()
	assert.Nil(t, errs)

	// register close-error component
	err = components.Register(closeErrorComponentName, &CloseErrorComponent{})
	assert.Nil(t, err)
	c = components.Get(regularComponentName)
	assert.Nil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.NotNil(t, c)
	errs = components.Close()
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.NotNil(t, errs[closeErrorComponentName])
	c = components.Get(closeErrorComponentName) // make sure close-init component is removed
	assert.Nil(t, c)

	// register all components
	err = components.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	err = components.Register(initErrorComponentName, &InitErrorComponent{})
	assert.NotNil(t, err)
	err = components.Register(closeErrorComponentName, &CloseErrorComponent{})
	assert.Nil(t, err)
	c = components.Get(regularComponentName)
	assert.NotNil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.NotNil(t, c)
	errs = components.Close()
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.NotNil(t, errs[closeErrorComponentName])
	c = components.Get(regularComponentName)
	assert.Nil(t, c)
	c = components.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = components.Get(closeErrorComponentName)
	assert.Nil(t, c)

	// Register() errors
	err = components.Register("", &RegularComponent{})
	assert.NotNil(t, err)
	err = components.Register(regularComponentName, nil)
	assert.NotNil(t, err)
	err = components.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	err = components.Register(regularComponentName, &RegularComponent{})
	assert.NotNil(t, err) // name already registered

	// Get() errors
	c = components.Get("")
	assert.Nil(t, c)
	c = components.Get(regularComponentName + "-2")
	assert.Nil(t, c)
}
