package server

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Server interface {
	Options() Options
	Init(...Option) error
	Start() error
	Stop() error
	String() string
}

var (
	DefaultAddress          = "localhost:9502"
	DefaultName             = "micro.server"
	DefaultVersion          = time.Now().Format("2006.01.02.15.04")
	DefaultId               = uuid.New().String()
	DefaultRegisterCheck    = func(context.Context) error { return nil }
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)

type Option func(*Options)

// NewServer returns a new server with options passed in
func NewServer(opt ...Option) Server {
	return newRpcServer(opt...)
}
