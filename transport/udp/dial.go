package udp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func Dial(address string) (io.ReadWriteCloser, error) {
	s, err := NewSocket(":0")
	if err != nil {
		return nil, fmt.Errorf("could not create client socket: %v", err)
	}

	peerAddr, err := ResolveAddr(address)
	if err != nil {
		return nil, fmt.Errorf("could not resolve peer address: %v", err)
	}

	synPacket := &Packet{
		Type:        SYN,
		Sequence:    rand.Uint32(),
		PeerAddress: binary.BigEndian.Uint32(peerAddr.IP.To4()),
		PeerPort:    uint16(peerAddr.Port),
		Payload:     []byte{},
	}

	err = s.Send(synPacket, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not send SYN packet: %v", err)
	}

	synAckPacket, err := s.Receive(10 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not receive SYN packet: %v", err)
	}

	if synAckPacket.Type != SYNACK {
		return nil, fmt.Errorf("could not synchronize with peer: incorrect SYN response type")
	}
	if synAckPacket.Sequence != synPacket.Sequence {
		return nil, fmt.Errorf("could not synchronize with peer: incorrect SYN response sequence")
	}

	ackPacket := &Packet{
		Type:        ACK,
		Sequence:    synPacket.Sequence,
		PeerAddress: synPacket.PeerAddress,
		PeerPort:    synPacket.PeerPort,
		Payload:     []byte{},
	}

	err = s.Send(ackPacket, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not send ACK packet: %v", err)
	}

	// TODO make connection instance
	return &Conn{}, nil
}
