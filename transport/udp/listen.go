package udp

import (
	"fmt"
	"log"

	"github.com/g-harel/http/transport/connection"
)

func Listen(port string) (*Listener, error) {
	log.SetPrefix("[SERVER]          ")
	log.SetFlags(0)
	log.Printf("Listen(%v)\n", port)

	s, err := NewSocket(port)
	if err != nil {
		return nil, fmt.Errorf("create sender socket: %v", err)
	}
	return &Listener{
		socket: s,
	}, nil
}

type Listener struct {
	socket *Socket
	server *Server
}

func (ln *Listener) Accept() (connection.Connection, error) {
	log.Printf("Listener.Accept()\n")

	s, err := NewServer(ln.socket)
	if err != nil {
		return nil, fmt.Errorf("create server: %v", err)
	}

	log.Println("connection established")

	ln.server = s

	return s, nil
}

func (ln *Listener) Close() error {
	return nil
}
