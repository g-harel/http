package tcp

import (
	"net"

	"github.com/g-harel/http/transport/connection"
)

var _ connection.Connection = &Connection{}

type Connection struct {
	net.Conn
}

func (c *Connection) Commit() error {
	return nil
}
