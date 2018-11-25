package udp

import (
	"fmt"
	"io"
	"time"
)

func Listen(port string) (*Listener, error) {
	s, err := NewSocket(port)
	if err != nil {
		return nil, fmt.Errorf("create sender socket: %v", err)
	}
	return &Listener{s}, nil
}

type Listener struct {
	socket *Socket
}

func (ln *Listener) Accept() (io.ReadWriteCloser, error) {
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

	err = ln.socket.Send(synAckPacket, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("send SYNACK packet: %v", err)
	}

	ackPacket, err := ln.socket.Receive(10 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("receive ACK packet: %v", err)
	}

	if ackPacket.Type != ACK {
		return nil, fmt.Errorf("synchronize with peer: incorrect ACK type")
	}
	if ackPacket.Sequence != synPacket.Sequence {
		return nil, fmt.Errorf("synchronize with peer: incorrect ACK response sequence")
	}

	// TODO make connection instance
	return &Conn{}, nil
}

func (ln *Listener) Close() error {
	return nil
}
