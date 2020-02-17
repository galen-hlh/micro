package server

import (
	"github.com/micro/go-micro/util/addr"
	"google.golang.org/grpc"
	"micro/gateway/metadata"
	"micro/gateway/mnet"
	"micro/gateway/registry"
	"micro/util/log"
	"net"
	"strings"
	"sync"
	"time"
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

	// used for first registration
	registered bool

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

	log.Logf("Micro [%s] Listening on %s", s.String(), ts.Addr().String())

	// use RegisterCheck func before register
	if err = s.opts.RegisterCheck(s.opts.Context); err != nil {
		log.Logf("Server %s-%s register check error: %s", config.Name, config.Id, err)
	} else {
		// announce self to the world
		if err = s.Register(); err != nil {
			log.Logf("Server %s-%s register error: %s", config.Name, config.Id, err)
		}
	}

	exit := make(chan bool)

	go func() {
		t := new(time.Ticker)

		// only process if it exists
		if s.opts.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(s.opts.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				s.RLock()
				registered := s.registered
				s.RUnlock()
				if err = s.opts.RegisterCheck(s.opts.Context); err != nil && registered {
					log.Logf("Server %s-%s register check error: %s, deregister it", config.Name, config.Id, err)
					// deregister self in case of error
					if err := s.Deregister(); err != nil {
						log.Logf("Server %s-%s deregister error: %s", config.Name, config.Id, err)
					}
				} else {
					if err := s.Register(); err != nil {
						log.Logf("Server %s-%s register error: %s", config.Name, config.Id, err)
					}
				}
			// wait for exit
			case ch = <-s.exit:
				t.Stop()
				close(exit)
				break Loop
			}
		}

		// deregister self
		if err := s.Deregister(); err != nil {
			log.Logf("Server %s-%s deregister error: %s", config.Name, config.Id, err)
		}

		s.Lock()
		swg := s.wg
		s.Unlock()

		// wait for requests to finish
		if swg != nil {
			swg.Wait()
		}

		// close transport listener
		ch <- ts.Close()
	}()

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
	return "grpc"
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

func (s *grpcServer) Register() error {
	var err error
	var advt, host, port string

	// parse address for host, port
	config := s.Options()

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	addr, err := addr.Extract(host)
	if err != nil {
		return err
	}

	// make copy of metadata
	md := make(metadata.Metadata)
	for k, v := range config.Metadata {
		md[k] = v
	}

	// mq-rpc(eg. nats) doesn't need the port. its addr is queue name.
	if port != "" {
		addr = mnet.HostPort(addr, port)
	}

	// register service
	node := &registry.Node{
		Id:       config.Name + "-" + config.Id,
		Address:  addr,
		Metadata: md,
	}

	node.Metadata["server"] = s.String()
	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["protocol"] = "mucp"

	s.RLock()

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
		//Endpoints: endpoints,
	}

	// get registered value
	registered := s.registered

	s.RUnlock()

	if !registered {
		log.Logf("Registry [%s] Registering node: %s", config.Registry.String(), node.Id)
	}

	// create registry options
	rOpts := []registry.RegisterOption{registry.RegisterTTL(config.RegisterTTL)}

	if err := config.Registry.Register(service, rOpts...); err != nil {
		return err
	}

	// already registered? don't need to register subscribers
	if registered {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	s.registered = true
	// set what we're advertising
	s.opts.Advertise = addr

	return nil
}

func (s *grpcServer) Deregister() error {
	var err error
	var advt, host, port string

	config := s.Options()

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	addr, err := addr.Extract(host)
	if err != nil {
		return err
	}

	// mq-rpc(eg. nats) doesn't need the port. its addr is queue name.
	if port != "" {
		addr = mnet.HostPort(addr, port)
	}

	node := &registry.Node{
		Id:      config.Name + "-" + config.Id,
		Address: addr,
	}

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	log.Logf("Registry [%s] Deregistering node: %s", config.Registry.String(), node.Id)
	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	s.Lock()

	if !s.registered {
		s.Unlock()
		return nil
	}

	s.registered = false

	s.Unlock()
	return nil
}

func newGrpcServer(opts ...Option) Server {
	options := newOptions(opts...)

	return &grpcServer{
		opts: options,
		exit: make(chan chan error),

		wg: wait(options.Context),
	}
}
