package transport

import (
	"net"
)

var _ Listener = &tcpListener{}

type tcpListener struct {
	ln net.Listener
}

func (ln *tcpListener) Accept() (Connection, error) {
	return ln.ln.Accept()
}

func (ln *tcpListener) Close() error {
	return ln.ln.Close()
}
