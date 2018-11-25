package tcp

import (
	"net"

	"github.com/g-harel/http/transport/connection"
)

func Listen(port string) (*Listener, error) {
	ln, err := net.Listen("tcp", port)
	return &Listener{ln}, err
}

type Listener struct {
	net.Listener
}

func (ln *Listener) Accept() (connection.Connection, error) {
	conn, err := ln.Listener.Accept()
	return &Connection{conn}, err
}

func (ln *Listener) Close() error {
	return ln.Listener.Close()
}
