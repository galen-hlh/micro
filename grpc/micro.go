package grpc

import (
	"micro/grpc/client"
	"micro/grpc/server"
)

type Service interface {
	// The grpc name
	Name() string
	//// Init initialises options
	Init(...Option)
	//// Options returns the current options
	Options() Options
	//// Client is used to call services
	Client() client.Client
	//// Server is for handling requests and events
	Server() server.Server
	//// Run the grpc
	Run() error
	//// The grpc implementation
	String() string
}

type Option func(*Options)

func NewService(opts ...Option) Service {
	return newService(opts...)
}
