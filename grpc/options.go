package grpc

import (
	"context"
	"micro/grpc/client"
	"micro/grpc/server"
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
	opt := Options{}

	return opt
}
