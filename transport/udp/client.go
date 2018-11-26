package udp

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/g-harel/http/transport/connection"
)

var _ connection.Connection = &Client{}

type Client struct {
	socket  *Socket
	packet  *Packet
	window  *bytes.Buffer
	mailbox *bytes.Buffer
}

func NewClient(socket *Socket, packet *Packet) (*Client, error) {
	client := &Client{
		socket:  socket,
		packet:  packet,
		window:  &bytes.Buffer{},
		mailbox: nil,
	}

	packet.Type = SYN

	err := socket.Send(packet, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send SYN packet: %v", err)
	}

	var synAckPacket *Packet
	select {
	case p := <-socket.Received:
		synAckPacket = p
	case <-time.After(Timeout):
		return nil, fmt.Errorf("wait for SYNACK: timeout")
	}

	if synAckPacket.Type != SYNACK {
		return nil, fmt.Errorf("synchronize with peer: incorrect SYN response type")
	}
	if synAckPacket.Sequence != packet.Sequence {
		return nil, fmt.Errorf("synchronize with peer: incorrect SYN response sequence")
	}

	ackPacket := &Packet{
		Type:        ACK,
		Sequence:    packet.Sequence,
		PeerAddress: packet.PeerAddress,
		PeerPort:    packet.PeerPort,
		Payload:     []byte{},
	}

	err = socket.Send(ackPacket, Timeout)
	if err != nil {
		return nil, fmt.Errorf("send ACK packet: %v", err)
	}

	return client, nil
}

func (c *Client) Read(b []byte) (int, error) {
	if c.mailbox == nil {
		return 0, io.EOF
	}

	loops := int(Timeout / Poll)
	for i := 0; i < loops; i++ {
		if c.mailbox.Len() != 0 {
			n, err := c.mailbox.Read(b)
			c.mailbox = nil
			return n, err
		}
		time.Sleep(Poll)
	}
	return 0, fmt.Errorf("read: timeout")
}

func (c *Client) Write(b []byte) (int, error) {
	return c.window.Write(b)
}

func (c *Client) Commit() error {
	c.packet.Sequence++
	c.mailbox = &bytes.Buffer{}

	p := &Packet{
		Sequence:    c.packet.Sequence,
		PeerAddress: c.packet.PeerAddress,
		PeerPort:    c.packet.PeerPort,
		Payload:     c.window.Bytes(),
	}

	err := c.socket.Send(p, Timeout)
	if err != nil {
		return fmt.Errorf("send packet: %v", err)
	}

	var ackPacket *Packet
	select {
	case p := <-c.socket.Received:
		ackPacket = p
	case <-time.After(Timeout):
		return fmt.Errorf("wait for ACK: timeout")
	}

	if ackPacket.Type != ACK {
		return fmt.Errorf("check ack packet: not ACK")
	}
	if ackPacket.Sequence != p.Sequence {
		return fmt.Errorf("check ack packet: sequence doesn't match")
	}

	c.window.Reset()

	_, err = c.mailbox.Write(ackPacket.Payload)
	if err != nil {
		return fmt.Errorf("write ack response: %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	return nil
}
