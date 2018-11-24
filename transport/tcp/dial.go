package tcp

import (
	"io"
	"net"
)

func Dial(address string) (io.ReadWriteCloser, error) {
	return net.Dial("tcp", address)
}
