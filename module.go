package module

//Option is the function type for initialization of the Modules
type Option func()

type moduler interface {
	Init(options ...Option) error
}

type Module struct {
	moduler
}

func (m *Module) ParseOpts(opts []Option) {
	for i := range opts {
		opts[i]()
	}
}
