package transport

import (
	"fmt"
	"net"
)

// Supported protocols.
const (
	TCP = "tcp"
	UDP = "udp"
)

// Connection represents a generic network connection.
type Connection interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

// Listener represents a generic network listener.
type Listener interface {
	Accept() (Connection, error)
	Close() error
}

func Dial(protocol string, address string) (Connection, error) {
	if protocol == TCP {
		return net.Dial(protocol, address)
	}

	return nil, fmt.Errorf("unrecognized protocol \"%v\"", protocol)
}

func Listen(protocol string, port string) (Listener, error) {
	if protocol == TCP {
		ln, err := net.Listen(protocol, port)
		if err != nil {
			return nil, fmt.Errorf("could not listen: %v", err)
		}
		return &tcpListener{ln}, nil
	}

	return nil, fmt.Errorf("unrecognized protocol \"%v\"", protocol)
}
