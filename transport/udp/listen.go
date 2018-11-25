package udp

import (
	"fmt"
	"log"

	"github.com/g-harel/http/transport/connection"
)

func Listen(port string) (*Listener, error) {
	log.SetPrefix("[SERVER] ")
	log.SetFlags(0)
	log.Printf("Listen(port: \"%v\")\n", port)

	s, err := NewSocket(port)
	if err != nil {
		return nil, fmt.Errorf("create sender socket: %v", err)
	}
	return &Listener{s}, nil
}

type Listener struct {
	socket *Socket
}

func (ln *Listener) Accept() (connection.Connection, error) {
	log.Printf("Accept()\n")

	synPacket, err := ln.socket.Receive(0)
	if err != nil {
		return nil, fmt.Errorf("receive SYN packet: %v", err)
	}

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

	err = ln.socket.Send(synAckPacket, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send SYNACK packet: %v", err)
	}

	ackPacket, err := ln.socket.Receive(Timeout)
	if err != nil {
		return nil, fmt.Errorf("receive ACK packet: %v", err)
	}

	if ackPacket.Type != ACK {
		return nil, fmt.Errorf("synchronize with peer: incorrect ACK type")
	}
	if ackPacket.Sequence != synPacket.Sequence {
		return nil, fmt.Errorf("synchronize with peer: incorrect ACK response sequence")
	}

	log.Printf("connection established\n")

	return NewConn(ln.socket, ackPacket.Sequence, ackPacket.PeerAddress, ackPacket.PeerPort), nil
}

func (ln *Listener) Close() error {
	return nil
}
