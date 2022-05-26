package mongoserve

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOpts *options.ClientOptions

// Option is a function that can be passed to Init
type Option interface {
	apply(opts *initOptions)
}

type initOptions struct {
	configer      configer.Configurable
	clientOptions *options.ClientOptions
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

// WithClientOpts set client options
func WithClientOpts(clientOpts *options.ClientOptions) Option {
	return OptFunc(func(opts *initOptions) {
		opts.clientOptions = clientOpts
	})
}
