package http

import (
	"net"
	"strings"
)

// Handler function handles the server's requests.
type Handler func(req *Request) (*Response, error)

func defaultHandler(req *Request) (*Response, error) {
	return &Response{
		Status:     "OK",
		StatusCode: 200,
		Body:       strings.NewReader("200 OK"),
	}, nil
}

// ErrorHandler should produce a response from an error.
type ErrorHandler func(err error) *Response

func defaultErrorHandler(err error) *Response {
	return &Response{
		StatusCode: 500,
		Status:     "Internal Server Error",
		Body:       strings.NewReader("500 Internal Server Error"),
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

	if s.Errors == nil {
		s.Errors = make(chan error)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case s.Errors <- err:
			}
			conn.Close()
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

		go handleConn(conn, *s)
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

func handleConn(conn net.Conn, s Server) {
	defer conn.Close()

	var res *Response

	req, err := ReadRequest(conn)
	if err != nil {
		res = s.errHandler(err)
	}

	res, err = s.handler(req)
	if err != nil {
		res = s.errHandler(err)
	}

	err = res.Fprint(conn)
	if err != nil {
		select {
		case s.Errors <- err:
		}
	}
}
