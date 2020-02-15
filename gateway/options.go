package gateway

import (
	"context"
	"micro/gateway/client"
	"micro/gateway/server"
	"time"
)

type Options struct {
	Client client.Client
	Server server.Server

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Server:  server.NewServer(),
		Context: context.Background(),
		Signal:  true,
	}

	return opt
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.Server.Init(server.Address(addr))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Init(server.Version(v))
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.Server.Init(server.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.Server.Init(server.RegisterInterval(t))
	}
}

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
