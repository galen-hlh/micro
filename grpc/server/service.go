package server

type Server interface {
	Options() Options
	Init(...Option) error
	Start() error
	Stop() error
	String() string
}

type Option func(*Options)
