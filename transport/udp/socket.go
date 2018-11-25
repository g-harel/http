package udp

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const Timeout = 2 * time.Second

func ResolveAddr(address string) (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, fmt.Errorf("resolve given address: %v", err)
	}
	if addr.IP == nil || addr.IP.IsUnspecified() {
		addr.IP = net.IP([]byte{127, 0, 0, 1})
	}

	return addr, nil
}

type Socket struct {
	Transport net.PacketConn
}

func NewSocket(address string) (*Socket, error) {
	addr, err := ResolveAddr(address)
	if err != nil {
		return nil, fmt.Errorf("resolve own address: %v", err)
	}

	conn, err := net.ListenPacket("udp4", addr.String())
	if err != nil {
		return nil, fmt.Errorf("create raw packet connection: %s", err)
	}

	return &Socket{conn}, nil
}

func (s *Socket) Send(p *Packet, timeout time.Duration) error {
	err := s.Transport.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return fmt.Errorf("set send timeout: %v", err)
	}

	// Overwritten address per assignment instructions (to use flaky router).
	addr, err := net.ResolveUDPAddr("udp4", ":"+os.Getenv("ROUTER_PORT"))
	if err != nil {
		return fmt.Errorf("resolve router address: %v", err)
	}

	_, err = s.Transport.WriteTo(p.Bytes(), addr)
	if err != nil {
		return fmt.Errorf("write packet: %v", err)
	}

	log.Printf("Socket.Send(%v, %v)", p.Type, p.Sequence)

	return nil
}

func (s *Socket) Receive(timeout time.Duration) (*Packet, error) {
	b := make([]byte, MaxPacketSize)

	// No deadline if timeout is zero.
	var deadline time.Time
	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}

	err := s.Transport.SetReadDeadline(deadline)
	if err != nil {
		return nil, fmt.Errorf("set timeout: %v", err)
	}

	n, _, err := s.Transport.ReadFrom(b)
	if err != nil {
		return nil, fmt.Errorf("read packet: %v", err)
	}

	p, err := (&Packet{}).Parse(b[:n])
	if err != nil {
		return nil, fmt.Errorf("parse packet: %v", err)
	}

	log.Printf("Socket.Receive(%v, %v)", p.Type, p.Sequence)

	return p, nil
}
