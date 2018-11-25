package udp

import (
	"fmt"
	"net"
	"os"
	"time"
)

func ResolveAddr(address string) (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, fmt.Errorf("could not resolve given address: %v", err)
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
		return nil, fmt.Errorf("could not resolve own address: %v", err)
	}

	conn, err := net.ListenPacket("udp4", addr.String())
	if err != nil {
		return nil, fmt.Errorf("could not create raw packet connection: %s", err)
	}

	return &Socket{conn}, nil
}

func (s *Socket) Send(p *Packet, timeout time.Duration) error {
	err := s.Transport.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return fmt.Errorf("could not set timeout: %v", err)
	}

	// Overwritten address per assignment instructions (to use flaky router).
	addr, err := net.ResolveUDPAddr("udp4", ":"+os.Getenv("ROUTER_PORT"))
	if err != nil {
		return fmt.Errorf("could not resolve router address: %v", err)
	}

	_, err = s.Transport.WriteTo(p.Bytes(), addr)
	if err != nil {
		return fmt.Errorf("could not write packet: %v", err)
	}

	return nil
}

func (s *Socket) Receive(timeout time.Duration) (*Packet, error) {
	b := make([]byte, MaxPacketSize)

	if timeout != 0 {
		err := s.Transport.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, fmt.Errorf("could not set timeout: %v", err)
		}
	}

	n, _, err := s.Transport.ReadFrom(b)
	if err != nil {
		return nil, fmt.Errorf("could not read packet: %v", err)
	}

	p, err := (&Packet{}).Parse(b[:n])
	if err != nil {
		return nil, fmt.Errorf("could not parse packet: %v", err)
	}

	return p, nil
}
