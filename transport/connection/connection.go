package connection

import (
	"io"
)

// Connection represents a generic network connection.
type Connection interface {
	io.ReadWriteCloser
	Commit() error
}
