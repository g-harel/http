package udp

import (
	"io"
)

func Listen(port string) (*Listener, error) {
	return &Listener{}, nil
}

type Listener struct {
}

func (ln *Listener) Accept() (io.ReadWriteCloser, error) {
	return nil, nil
}

func (ln *Listener) Close() error {
	return nil
}
