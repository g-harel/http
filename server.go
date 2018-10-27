package http

import (
	"fmt"
	"net"
	"strings"
)

// Handler function handles the server's requests.
type Handler func(req *Request) (*Response, error)

func defaultHandler(req *Request) (*Response, error) {
	return &Response{
		Status:     "OK",
		StatusCode: 200,
		Body:       strings.NewReader("200 OK\n"),
	}, nil
}

// ErrorHandler should produce a response from an error.
type ErrorHandler func(err error) *Response

func defaultErrorHandler(err error) *Response {
	return &Response{
		StatusCode: 500,
		Status:     "Internal Server Error",
		Body:       strings.NewReader("500 Internal Server Error\n"),
	}
}

// Server is used to respond to http requests.
type Server struct {
	handler    Handler
	errHandler ErrorHandler

	// Errors that cannot be handled by errorHandler are sent to this channel.
	Errors chan error
}

// Listen listens for incoming requests on the requested port.
func (s *Server) Listen(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			conn.Close()
			s.throw(err)
			continue
		}

		// Server config is copied before being sent to a separate thread.
		config := Server{
			Errors:     s.Errors,
			handler:    s.handler,
			errHandler: s.errHandler,
		}
		if config.errHandler == nil {
			config.errHandler = defaultErrorHandler
		}
		if config.handler == nil {
			config.handler = defaultHandler
		}

		go handleConn(conn, config)
	}
}

// Use configures the server to use the given Handler.
func (s *Server) Use(h Handler) {
	s.handler = h
}

// Err configures the server to use the given ErrorHandler.
func (s *Server) Err(h ErrorHandler) {
	s.errHandler = h
}

// Non-blocking send to the server's error channel.
func (s *Server) throw(err error) {
	if s.Errors != nil {
		select {
		case s.Errors <- fmt.Errorf("ERROR: %v", err):
		default:
		}
	}
}

func handleConn(conn net.Conn, s Server) {
	var res *Response
	defer func() {
		err := res.Fprint(conn)
		if err != nil {
			s.throw(err)
		}
		conn.Close()
	}()

	req, err := ReadRequest(conn)
	if err != nil {
		s.throw(err)
		res = s.errHandler(err)
		return
	}

	res, err = s.handler(req)
	if err != nil {
		s.throw(err)
		res = s.errHandler(err)
	}
}
