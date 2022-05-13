package mongoserve

import "github.com/jacksonCLyu/ridi-faces/pkg/configer"

// Option is a function that can be passed to Init
type Option interface {
	apply(opts *initOptions)
}

type initOptions struct {
	configer configer.Configurable
}

type OptFunc func(opts *initOptions)

func (f OptFunc) apply(opts *initOptions) {
	f(opts)
}

// WithConfig set config
func WithConfig(config configer.Configurable) Option {
	return OptFunc(func(opts *initOptions) {
		opts.configer = config
	})
}
