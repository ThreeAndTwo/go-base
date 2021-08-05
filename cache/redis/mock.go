package redis

import (
	"github.com/alicebob/miniredis/v2"
)

var _ Interface = &Mocker{}

type Mocker struct {
	server *miniredis.Miniredis
	*Redis
}

func NewMocker() (*Mocker, error) {
	server, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	client, err := New(&Config{Addr: server.Addr()})
	if err != nil {
		return nil, err
	}
	return &Mocker{
		server: server,
		Redis:  client,
	}, nil
}
