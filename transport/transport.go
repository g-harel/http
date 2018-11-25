package transport

import (
	"fmt"
	"math/rand"

	"github.com/g-harel/http/transport/connection"
	"github.com/g-harel/http/transport/tcp"
	"github.com/g-harel/http/transport/udp"
)

// Supported protocols.
const (
	TCP = "tcp"
	UDP = "udp"
)

// EOF is a special sequence to signal the end of a message (very sketchy).
var EOF = []byte{byte(rand.Uint64()), byte(rand.Uint64()), byte(rand.Uint64()), byte(rand.Uint64()), byte(rand.Uint64())}

// Listener represents a generic network listener.
type Listener interface {
	Accept() (connection.Connection, error)
	Close() error
}

func Dial(protocol string, address string) (connection.Connection, error) {
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
