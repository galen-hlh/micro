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

type rpcServer struct {
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

func (s *rpcServer) Options() Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

func (s *rpcServer) Start() error {
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

func (s *rpcServer) Stop() error {
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

func (s *rpcServer) Init(opts ...Option) error {
	s.configure(opts...)
	return nil
}

func (s *rpcServer) String() string {
	return "rpc"
}

func (s *rpcServer) configure(opts ...Option) {
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

func (s *rpcServer) getMaxMsgSize() int {
	if s.opts.maxMsgSize == 0 {
		return DefaultMaxMsgSize
	} else {
		return s.opts.maxMsgSize
	}
}

func newRpcServer(opts ...Option) Server {
	options := newOptions(opts...)

	return &rpcServer{
		opts: options,
		exit: make(chan chan error),

		wg: wait(options.Context),
	}
}
