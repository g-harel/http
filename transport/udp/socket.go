package udp

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const Timeout = 2 * time.Second
const Poll = 100 * time.Millisecond

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
	Received  chan *Packet
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

	s := &Socket{
		Transport: conn,
		Received:  make(chan *Packet),
	}

	go func() {
		buffer := make([]byte, MaxPacketSize)
		for {
			n, _, err := s.Transport.ReadFrom(buffer)
			if err != nil {
				log.Printf("error: read packet: %v", err)
				continue
			}

			p, err := (&Packet{}).Parse(buffer[:n])
			if err != nil {
				log.Printf("parse packet: %v", err)
				continue
			}

			log.Printf("Socket.Receive(%v, %v)", p.Type, p.Sequence)

			s.Received <- p
		}
	}()

	return s, nil
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
