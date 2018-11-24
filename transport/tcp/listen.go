package tcp

import (
	"io"
	"net"
)

func Listen(port string) (*Listener, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	return &Listener{ln}, nil
}

type Listener struct {
	net.Listener
}

func (ln *Listener) Accept() (io.ReadWriteCloser, error) {
	return ln.Listener.Accept()
}

func (ln *Listener) Close() error {
	return ln.Listener.Close()
}
