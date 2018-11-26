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
	log.Printf("Dial(%v)\n", address)

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
		Sequence:    rand.Uint32() % 128,
		PeerAddress: binary.BigEndian.Uint32(peerAddr.IP.To4()),
		PeerPort:    uint16(peerAddr.Port),
		Payload:     []byte{},
	}

	client, err := NewClient(s, synPacket)
	if err != nil {
		return nil, fmt.Errorf("create client: %v", err)
	}

	log.Println("connection established")

	return client, nil
}
