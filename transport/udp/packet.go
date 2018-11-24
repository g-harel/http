package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Valid packet types.
const (
	ACK uint8 = iota + 1
	SYN
	SYNACK
	NAK
)

type Packet struct {
	Type        uint8
	Sequence    uint32
	PeerAddress uint32
	PeerPort    uint16
	Payload     []byte
}

func (p *Packet) Bytes() []byte {
	b := bytes.Buffer{}

	b.WriteByte(p.Type)
	binary.Write(&b, binary.BigEndian, p.Sequence)
	binary.Write(&b, binary.BigEndian, p.PeerAddress)
	binary.Write(&b, binary.BigEndian, p.PeerPort)
	binary.Write(&b, binary.BigEndian, p.Payload)

	return b.Bytes()
}

func (p *Packet) Parse(data []byte) (*Packet, error) {
	if len(data) < 11 {
		return nil, fmt.Errorf("missing packet header data")
	}
	if len(data) > 1024 {
		return nil, fmt.Errorf("packet is too large")
	}

	p.Type = data[0]
	p.Sequence = binary.BigEndian.Uint32(data[1:5])
	p.PeerAddress = binary.BigEndian.Uint32(data[5:9])
	p.PeerPort = binary.BigEndian.Uint16(data[9:11])
	p.Payload = data[11:]

	return p, nil
}
