package server

import "context"

type Options struct {
	Name      string
	Address   string
	Advertise string
	Id        string
	Version   string

	Context context.Context
}
