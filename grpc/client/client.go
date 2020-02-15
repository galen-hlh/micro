package client

type Client interface {
	Init(...Option) error
	Options() Options
}

type Option func(*Options)
