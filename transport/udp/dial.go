package udp

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"

	"github.com/g-harel/http/transport/connection"
)

func Dial(address string) (connection.Connection, error) {
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(0)
	log.Printf("Dial(address: \"%v\")\n", address)

	s, err := NewSocket(":0")
	if err != nil {
		return nil, fmt.Errorf("create client socket: %v", err)
	}

	peerAddr, err := ResolveAddr(address)
	if err != nil {
		return nil, fmt.Errorf("resolve peer address: %v", err)
	}

	synPacket := &Packet{
		Type:        SYN,
		Sequence:    rand.Uint32(),
		PeerAddress: binary.BigEndian.Uint32(peerAddr.IP.To4()),
		PeerPort:    uint16(peerAddr.Port),
		Payload:     []byte{},
	}

	err = s.Send(synPacket, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send SYN packet: %v", err)
	}

	synAckPacket, err := s.Receive(Timeout)
	if err != nil {
		return nil, fmt.Errorf("receive SYN packet: %v", err)
	}

	if synAckPacket.Type != SYNACK {
		return nil, fmt.Errorf("synchronize with peer: incorrect SYN response type")
	}
	if synAckPacket.Sequence != synPacket.Sequence {
		return nil, fmt.Errorf("synchronize with peer: incorrect SYN response sequence")
	}

	ackPacket := &Packet{
		Type:        ACK,
		Sequence:    synPacket.Sequence,
		PeerAddress: synPacket.PeerAddress,
		PeerPort:    synPacket.PeerPort,
		Payload:     []byte{},
	}

	err = s.Send(ackPacket, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send ACK packet: %v", err)
	}

	log.Printf("connection established\n")

	return NewConn(s, ackPacket.Sequence, ackPacket.PeerAddress, ackPacket.PeerPort), nil
}
