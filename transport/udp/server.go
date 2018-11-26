package udp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/g-harel/http/transport/connection"
)

var _ connection.Connection = &Client{}

type Server struct {
	socket  *Socket
	window  []*Packet
	mailbox *bytes.Buffer
	packet  *Packet
	close   chan bool
}

func NewServer(socket *Socket) (*Server, error) {
	s := &Server{
		socket:  socket,
		window:  []*Packet{},
		mailbox: &bytes.Buffer{},
		close:   make(chan bool),
	}

	synPacket := <-socket.Received

	if synPacket.Type != SYN {
		return nil, fmt.Errorf("synchronize with peer: incorrect SYN type")
	}

	synAckPacket := &Packet{
		Type:        SYNACK,
		Sequence:    synPacket.Sequence,
		PeerAddress: synPacket.PeerAddress,
		PeerPort:    synPacket.PeerPort,
		Payload:     []byte{},
	}

	err := socket.Send(synAckPacket, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send SYNACK packet: %v", err)
	}

	go func() {
		for {
			select {
			case p := <-socket.Received:
				s.window = append(s.window, p)
			case <-s.close:
				return
			}
		}
	}()

	loops := int(Timeout / Poll)
	for i := 0; i < loops; i++ {
		for i, packet := range s.window {
			if packet.Type != ACK {
				continue
			}

			s.window = append(s.window[:i], s.window[i+1:]...)

			if packet.Sequence != synPacket.Sequence {
				return nil, fmt.Errorf("synchronize with peer: incorrect ACK response sequence")
			}

			return s, nil
		}
		time.Sleep(Poll)
	}

	return nil, fmt.Errorf("connection handshake: ack timeout")
}

func (s *Server) Read(b []byte) (int, error) {
	loops := int(Timeout / Poll)
	for i := 0; i < loops; i++ {
		if len(s.window) != 0 {
			packet := s.window[0]
			s.window = s.window[1:]
			s.packet = packet

			if len(b) < len(packet.Payload) {
				return 0, fmt.Errorf("read packet: read buffer too small")
			}

			copy(b, packet.Payload)

			return len(packet.Payload), io.EOF
		}
		time.Sleep(Poll)
	}

	return 0, fmt.Errorf("read: timeout")
}

func (s *Server) Write(b []byte) (int, error) {
	return s.mailbox.Write(b)
}

func (s *Server) Commit() error {
	p := &Packet{
		Type:        ACK,
		Sequence:    s.packet.Sequence,
		PeerAddress: s.packet.PeerAddress,
		PeerPort:    s.packet.PeerPort,
		Payload:     s.mailbox.Bytes(),
	}

	err := s.socket.Send(p, Timeout)
	if err != nil {
		return fmt.Errorf("send packet: %v", err)
	}

	s.mailbox.Reset()

	return nil
}

func (s *Server) Close() error {
	log.Println("Server.Close()")
	s.close <- true
	return nil
}
