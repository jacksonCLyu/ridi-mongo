package mongoserve

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultClientOpts *options.ClientOptions

// GetLOptions returns default global client options
func GetLOptions() *options.ClientOptions {
	return defaultClientOpts
}

// SetLOptions set default global client options
func SetLOptions(clientOpts *options.ClientOptions) {
	defaultClientOpts = clientOpts
}

// Option is a function that can be passed to Init
type Option interface {
	apply(opts *initOptions)
}

type initOptions struct {
	serveName     string
	clientOptions *options.ClientOptions
}

type OptFunc func(opts *initOptions)

func (f OptFunc) apply(opts *initOptions) {
	f(opts)
}

// WithClientOpts set client options
func WithClientOpts(clientOpts *options.ClientOptions) Option {
	return OptFunc(func(opts *initOptions) {
		opts.clientOptions = clientOpts
	})
}
