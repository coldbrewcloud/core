package components

var systemComponents ComponentManager

func Register(name string, component Component) error {
	return systemComponents.Register(name, component)
}

func Get(name string) Component {
	return systemComponents.Get(name)
}

func Close() map[string]error {
	return systemComponents.Close()
}
