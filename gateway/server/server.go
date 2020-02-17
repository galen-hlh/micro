package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"
)

type Server interface {
	Options() Options
	Init(...Option) error
	Start() error
	Stop() error
	String() string
	GetGrpcServer() *grpc.Server
}

var (
	DefaultAddress          = "localhost:9501"
	DefaultName             = "micro.server"
	DefaultRegisterAddress  = "localhost:2791"
	DefaultVersion          = time.Now().Format("2006.01.02.15.04")
	DefaultId               = uuid.New().String()
	DefaultRegisterCheck    = func(context.Context) error { return nil }
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)

type Option func(*Options)

// NewServer returns a new server with options passed in
func NewServer(opt ...Option) Server {
	return newGrpcServer(opt...)
}
