package components

type Component interface {
	Init() error
	Close() error
}
