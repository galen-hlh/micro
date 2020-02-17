package server

import (
	"google.golang.org/grpc"
	"micro/util/log"
	"net"
	"sync"
)

var (
	// DefaultMaxMsgSize define maximum message size that server can send
	// or receive.  Default value is 4MB.
	DefaultMaxMsgSize = 1024 * 1024 * 4
)

type grpcServer struct {
	sync.RWMutex

	//grpc server
	srv *grpc.Server

	exit chan chan error
	opts Options

	// marks the serve as started
	started bool
	// graceful exit
	wg *sync.WaitGroup
}

func (s *grpcServer) Options() Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

func (s *grpcServer) Start() error {
	s.RLock()
	if s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	config := s.Options()

	// start listening on the transport
	ts, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}

	log.Logf("Transport [%s] Listening on %s", s.String(), ts.Addr().String())

	// micro: go ts.Accept(s.accept)
	go func() {
		if err := s.srv.Serve(ts); err != nil {
			log.Log("gRPC Server start error: ", err)
		}
	}()

	// mark the server as started
	s.Lock()
	s.started = true
	s.Unlock()

	return nil
}

func (s *grpcServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	ch := make(chan error)
	s.exit <- ch

	err := <-ch
	s.Lock()
	s.started = false
	s.Unlock()

	return err
}

func (s *grpcServer) Init(opts ...Option) error {
	s.configure(opts...)
	return nil
}

func (s *grpcServer) String() string {
	return "rpc"
}

func (s *grpcServer) configure(opts ...Option) {
	// Don't reprocess where there's no config
	if len(opts) == 0 && s.srv != nil {
		return
	}

	for _, o := range opts {
		o(&s.opts)
	}

	maxMsgSize := s.getMaxMsgSize()

	gopts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}

	s.srv = grpc.NewServer(gopts...)
}

func (s *grpcServer) GetGrpcServer() *grpc.Server {
	return s.srv
}

func (s *grpcServer) getMaxMsgSize() int {
	if s.opts.Context == nil {
		return DefaultMaxMsgSize
	}
	g, ok := s.opts.Context.Value(maxMsgSizeKey{}).(int)
	if !ok {
		return DefaultMaxMsgSize
	}
	return g
}

func newGrpcServer(opts ...Option) Server {
	options := newOptions(opts...)

	return &grpcServer{
		opts: options,
		exit: make(chan chan error),

		wg: wait(options.Context),
	}
}
