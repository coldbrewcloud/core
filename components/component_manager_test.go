package components_test

import (
	"fmt"
	"testing"

	"github.com/coldbrewcloud/core/components"
	"github.com/stretchr/testify/assert"
)

var (
	regularComponentName    = "c1"
	initErrorComponentName  = "c2"
	closeErrorComponentName = "c3"
)

type RegularComponent struct{}

func (c *RegularComponent) Init() error  { return nil }
func (c *RegularComponent) Close() error { return nil }

type InitErrorComponent struct{}

func (c *InitErrorComponent) Init() error  { return fmt.Errorf("Init error") }
func (c *InitErrorComponent) Close() error { return nil }

type CloseErrorComponent struct{}

func (c *CloseErrorComponent) Init() error  { return nil }
func (c *CloseErrorComponent) Close() error { return fmt.Errorf("Close error") }

func TestComponentManager(t *testing.T) {
	manager := &components.ComponentManager{}

	// initial state
	c := manager.Get(regularComponentName)
	assert.Nil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs := manager.Close()
	assert.Nil(t, errs)

	// register regular component
	err := manager.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	c = manager.Get(regularComponentName)
	assert.NotNil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs = manager.Close()
	assert.Nil(t, errs)

	// register init-error component
	err = manager.Register(initErrorComponentName, &InitErrorComponent{})
	assert.NotNil(t, err)
	c = manager.Get(regularComponentName)
	assert.Nil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.Nil(t, c)
	errs = manager.Close()
	assert.Nil(t, errs)

	// register close-error component
	err = manager.Register(closeErrorComponentName, &CloseErrorComponent{})
	assert.Nil(t, err)
	c = manager.Get(regularComponentName)
	assert.Nil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.NotNil(t, c)
	errs = manager.Close()
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.NotNil(t, errs[closeErrorComponentName])
	c = manager.Get(closeErrorComponentName) // make sure close-init component is removed
	assert.Nil(t, c)

	// register all components
	err = manager.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	err = manager.Register(initErrorComponentName, &InitErrorComponent{})
	assert.NotNil(t, err)
	err = manager.Register(closeErrorComponentName, &CloseErrorComponent{})
	assert.Nil(t, err)
	c = manager.Get(regularComponentName)
	assert.NotNil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.NotNil(t, c)
	errs = manager.Close()
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.NotNil(t, errs[closeErrorComponentName])
	c = manager.Get(regularComponentName)
	assert.Nil(t, c)
	c = manager.Get(initErrorComponentName)
	assert.Nil(t, c)
	c = manager.Get(closeErrorComponentName)
	assert.Nil(t, c)

	// Register() errors
	err = manager.Register("", &RegularComponent{})
	assert.NotNil(t, err)
	err = manager.Register(regularComponentName, nil)
	assert.NotNil(t, err)
	err = manager.Register(regularComponentName, &RegularComponent{})
	assert.Nil(t, err)
	err = manager.Register(regularComponentName, &RegularComponent{})
	assert.NotNil(t, err) // name already registered

	// Get() errors
	c = manager.Get("")
	assert.Nil(t, c)
	c = manager.Get(regularComponentName + "-2")
	assert.Nil(t, c)
}

type ValueComponent struct {
	Value int
	Dest  *[]int
}

func (c *ValueComponent) Init() error  { return nil }
func (c *ValueComponent) Close() error { *c.Dest = append(*c.Dest, c.Value); return nil }

func TestComponentManagerCloseOrder(t *testing.T) {
	manager := &components.ComponentManager{}
	dest := &[]int{}

	err := manager.Register("c1", &ValueComponent{Value: 1, Dest: dest})
	assert.Nil(t, err)
	err = manager.Register("c2", &ValueComponent{Value: 2, Dest: dest})
	assert.Nil(t, err)
	err = manager.Register("c3", &ValueComponent{Value: 3, Dest: dest})
	assert.Nil(t, err)
	errs := manager.Close() // components should be closed in reverse order
	assert.Nil(t, errs)

	assert.Equal(t, []int{3, 2, 1}, *dest)
}
