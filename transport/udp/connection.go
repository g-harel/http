package udp

import (
	"fmt"
	"io"
	"net"
	"os"
)

var RouterPort = os.Getenv("ROUTER_PORT")

var _ io.ReadWriteCloser = &Conn{}

type Conn struct {
	To  *net.UDPAddr
	Raw net.PacketConn

	Chan chan Packet
	Errs chan error

	// sendWindow [][]byte
	// sendChan   chan Packet

	// recvWindow [][]byte
	// recvChan   chan Packet

	// Sequence    uint32
	// PeerAddress uint32
	// PeerPort    uint16
}

func NewConn(from string) (*Conn, error) {
	c, err := net.ListenPacket("udp4", from)
	if err != nil {
		return nil, fmt.Errorf("could not create raw packet connection: %s", err)
	}

	conn := &Conn{
		Raw:  c,
		Chan: make(chan Packet),
		Errs: make(chan error),
	}

	return conn, nil
}

func (c *Conn) Peer(addr string) error {
	to, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return fmt.Errorf("could not resolve peer address: %v", err)
	}
	c.To = to
	return nil
}

func (c *Conn) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func (c *Conn) Write(b []byte) (int, error) {
	return len(b), nil
}

func (c *Conn) Close() error {
	// TODO connection close handshake
	if c.Raw != nil {
		return c.Raw.Close()
	}
	return nil
}
