package http

import (
	"fmt"
	"log"

	"github.com/g-harel/http/transport"
	"github.com/g-harel/http/transport/connection"
)

// Handler function handles the server's requests.
type Handler func(req *Request) (*Response, error)

// ErrorHandler should produce a response from an error.
type ErrorHandler func(err error) *Response

// Server is used to respond to http requests.
type Server struct {
	handler    Handler
	errHandler ErrorHandler

	// Errors that cannot be handled by errHandler are sent to this channel.
	ErrChan chan error
}

// Listen listens for incoming requests on the requested port.
func (s *Server) Listen(port string) error {
	ln, err := transport.Listen(transportProtocol, port)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			s.throw(fmt.Errorf("accept connection: %v", err))
			continue
		}

		config := Server{
			ErrChan:    s.ErrChan,
			handler:    s.handler,
			errHandler: s.errHandler,
		}
		if config.errHandler == nil {
			config.errHandler = defaultErrorHandler
		}
		if config.handler == nil {
			config.handler = defaultHandler
		}

		// Synchronously handle request.
		handleConn(conn, config)
	}
}

// Use configures the server to use the given Handler.
func (s *Server) Use(h Handler) {
	s.handler = h
}

// Catch configures the server to use the given ErrorHandler.
func (s *Server) Catch(h ErrorHandler) {
	s.errHandler = h
}

// Sends error to the server's error channel.
// It is assumed that the values will be received if channel is non-nil.
func (s *Server) throw(err error) {
	if s.ErrChan != nil {
		s.ErrChan <- fmt.Errorf("server: %v", err)
	}
}

func defaultHandler(req *Request) (*Response, error) {
	return NewResponse(200), nil
}

func defaultErrorHandler(err error) *Response {
	return NewResponse(500)
}

func handleConn(conn connection.Connection, s Server) {
	var res *Response
	defer func() {
		err := res.Fprint(conn)
		if err != nil {
			s.throw(fmt.Errorf("write response: %v", err))
		}

		err = conn.Commit()
		if err != nil {
			s.throw(fmt.Errorf("commit response: %v", err))
		}

		conn.Close()
	}()

	req, err := ReadRequest(conn)
	if err != nil {
		s.throw(fmt.Errorf("read request: %v", err))
		res = s.errHandler(err)
		return
	}

	log.Printf("handl(req)\n")
	res, err = s.handler(req)
	if err != nil {
		s.throw(fmt.Errorf("handle request: %v", err))
		res = s.errHandler(err)
	}
}
