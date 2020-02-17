package server

import (
	"context"
	"micro/gateway/registry"
	"time"
)

type maxMsgSizeKey struct{}

type Options struct {
	Metadata  map[string]string
	Name      string
	Address   string
	Advertise string
	Id        string
	Version   string

	// Registry service discovery
	Registry registry.Registry

	// service register address
	RegisterAddress []string

	// max msg size
	maxMsgSize int

	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// RegisterCheck run func before registry service
func RegisterCheck(fn func(context.Context) error) Option {
	return func(o *Options) {
		o.RegisterCheck = fn
	}
}

// Register the service with a TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// Register the service with at interval
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

func newOptions(opt ...Option) Options {
	opts := Options{
		Metadata:         map[string]string{},
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
	}

	for _, o := range opt {
		o(&opts)
	}

	//if opts.Registry == nil {
	//	opts.Registry = registry.DefaultRegistry
	//}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = DefaultRegisterCheck
	}

	if len(opts.Address) == 0 {
		opts.Address = DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = DefaultName
	}

	if len(opts.Id) == 0 {
		opts.Id = DefaultId
	}

	if len(opts.Version) == 0 {
		opts.Version = DefaultVersion
	}

	if len(opts.RegisterAddress) == 0 {
		opts.RegisterAddress = []string{DefaultRegisterAddress}
	}

	return opts
}
