package tcp

import (
	"net"

	"github.com/g-harel/http/transport/connection"
)

func Dial(address string) (connection.Connection, error) {
	conn, err := net.Dial("tcp", address)
	return &Connection{conn}, err
}
