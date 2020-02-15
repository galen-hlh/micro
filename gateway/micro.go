package gateway

import (
	"micro/gateway/client"
	"micro/gateway/server"
)

type Gateway interface {
	// The application name
	Name() string
	//// Init initialises options
	Init(...Option)
	//// Options returns the current options
	Options() Options
	//// Client is used to call services
	Client() client.Client
	//// Server is for handling requests and events
	Server() server.Server
	//// Run the application
	Run() error
	//// The application implementation
	String() string
}

type Option func(*Options)

func NewService(opts ...Option) Gateway {
	return newService(opts...)
}
