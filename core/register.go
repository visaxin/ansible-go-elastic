package core

type Context interface{}

type Do func(ctx Context, target []byte) error

func (this Do) Do(ctx Context, target []byte) error {
	return this.Do(ctx, target)
}

var register func(id string) MetadataRegister

// one way to register a cluster must implements the interface
type MetadataRegister interface {
	Do(ctx Context, metaData []byte) error
}

func Register(fn func() MetadataRegister) MetadataRegister {
	if register == nil {
		panic("cannot register metadata saver more than once")
	}
	register = func(_ string) MetadataRegister {
		return fn()
	}
	return fn()
}
