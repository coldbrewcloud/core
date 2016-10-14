package components

import (
	"fmt"
	"sync"

	"github.com/coldbrewcloud/coldbrew-core/errors"
)

type namedComponent struct {
	Name      string
	Component Component
}

type ComponentManager struct {
	components []namedComponent
	lock       sync.Mutex
}

func (cm *ComponentManager) Register(name string, component Component) error {
	if name == "" {
		return errors.EmptyParamError("name")
	}
	if component == nil {
		return errors.NilParamError("component")
	}

	cm.lock.Lock()
	defer cm.lock.Unlock()

	if cm.components == nil {
		cm.components = []namedComponent{}
	}

	if cm.Get(name) != nil {
		return fmt.Errorf("Component %q already registered", name)
	}

	if err := component.Init(); err != nil {
		return fmt.Errorf("Component %q init error: %s", name, err.Error())
	}

	cm.components = append(cm.components, namedComponent{
		Name:      name,
		Component: component,
	})

	return nil
}

func (cm *ComponentManager) Get(name string) Component {
	if name == "" {
		return nil
	}

	for _, c := range cm.components {
		if c.Name == name {
			return c.Component
		}
	}

	return nil
}

func (cm *ComponentManager) Close() map[string]error {
	errs := make(map[string]error)

	cm.lock.Lock()
	defer cm.lock.Unlock()

	// close in reverse order
	for i := len(cm.components) - 1; i >= 0; i-- {
		c := cm.components[i]

		if err := c.Component.Close(); err != nil {
			errs[c.Name] = err
		}
	}

	cm.components = []namedComponent{}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
