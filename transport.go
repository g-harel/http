package http

import (
	"fmt"
	"os"

	"github.com/g-harel/http/transport"
)

var transportProtocol string

func init() {
	protocol := os.Getenv("TRANSPORT_PROTOCOL")
	transportProtocol = protocol

	if protocol == transport.TCP {
		return
	}
	if protocol == transport.UDP {
		return
	}
	if protocol == "" {
		transportProtocol = transport.TCP
		return
	}

	panic(fmt.Errorf("unknown transport protocol"))
}
