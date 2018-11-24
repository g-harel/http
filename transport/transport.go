package transport

import (
	"fmt"
	"io"

	"github.com/g-harel/http/transport/tcp"
	"github.com/g-harel/http/transport/udp"
)

// Supported protocols.
const (
	TCP = "tcp"
	UDP = "udp"
)

// Connection represents a generic network connection.
type Connection = io.ReadWriteCloser

// Listener represents a generic network listener.
type Listener interface {
	Accept() (Connection, error)
	Close() error
}

func Dial(protocol string, address string) (Connection, error) {
	if protocol == TCP {
		return tcp.Dial(address)
	}
	if protocol == UDP {
		return udp.Dial(address)
	}

	return nil, fmt.Errorf("unrecognized protocol \"%v\"", protocol)
}

func Listen(protocol string, port string) (Listener, error) {
	if protocol == TCP {
		return tcp.Listen(port)
	}
	if protocol == UDP {
		return udp.Listen(port)
	}

	return nil, fmt.Errorf("unrecognized protocol \"%v\"", protocol)
}
