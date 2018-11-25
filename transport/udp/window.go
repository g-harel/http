package udp

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Window struct {
	buffer     map[uint32]*Packet
	lock       sync.Mutex
	sequence   uint32
	socket     *Socket
	errChan    chan error
	swallowAck bool
}

func NewWindow(socket *Socket, sequence uint32) *Window {
	w := &Window{
		buffer:   map[uint32]*Packet{},
		lock:     sync.Mutex{},
		sequence: sequence,
		socket:   socket,
	}

	go func() {
		for {
			p, err := socket.Receive(0)
			if err != nil {
				w.throw(fmt.Errorf("receive packet: %v", err))
			}

			if w.swallowAck && p.Type == ACK {
				continue
			}

			w.lock.Lock()
			log.Printf("Window.Buffer(%v, %v)", p.Type, p.Sequence)
			w.buffer[p.Sequence] = p
			w.lock.Unlock()

			if p.Type != 0 {
				continue
			}

			ackPacket := &Packet{
				Type:        ACK,
				Sequence:    p.Sequence,
				PeerAddress: p.PeerAddress,
				PeerPort:    p.PeerPort,
				Payload:     []byte{},
			}
			err = socket.Send(ackPacket, Timeout)
			if err != nil {
				w.throw(fmt.Errorf("ack packet: %v", err))
			}
		}
	}()

	return w
}

func (w *Window) Read(timeout time.Duration) (*Packet, error) {
	poll := 400 * time.Millisecond
	loops := int(1 + timeout/poll)

	for i := 0; i < loops; i++ {
		s := ""
		for k, _ := range w.buffer {
			s += fmt.Sprintf("%v ", k)
		}
		log.Printf("Window.Poll(%v) in [%v]", w.sequence, s)

		w.lock.Lock()
		expectedPacket := w.buffer[w.sequence]
		if expectedPacket != nil {
			log.Printf("Window.Read(%v)", w.sequence)
			delete(w.buffer, w.sequence)
			w.sequence++
			w.lock.Unlock()
			return expectedPacket, nil
		}
		w.lock.Unlock()
		time.Sleep(poll)
	}

	return nil, fmt.Errorf("window: read timeout")
}

// Errors are blocking if the error channel is defined.
func (w *Window) throw(err error) {
	if w.errChan != nil {
		w.errChan <- fmt.Errorf("window: %v", err)
	}
}
