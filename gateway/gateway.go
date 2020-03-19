package gateway

import (
	"micro/gateway/client"
	"micro/gateway/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type gateway struct {
	opts Options

	once sync.Once
}

func (s *gateway) Name() string {
	return s.opts.Server.Options().Name
}

func (s *gateway) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		s.opts.Server.Init()
	})
}

func (s *gateway) Options() Options {
	return s.opts
}

func (s *gateway) Client() client.Client {
	return s.opts.Client
}

func (s *gateway) Server() server.Server {
	return s.opts.Server
}

func (s *gateway) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *gateway) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *gateway) Run() error {

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

func (s *gateway) String() string {
	return "micro gateway"
}

func newService(opts ...Option) Gateway {
	options := newOptions(opts...)

	return &gateway{
		opts: options,
	}
}
