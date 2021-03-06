package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Response represents a received HTTP response.
type Response struct {
	Version    string
	Status     string
	StatusCode int
	Headers    *Headers
	Body       io.Reader

	conn net.Conn
}

// NewResponse creates a response for the given status and message.
// If no strings are given to fill the body, the status message is used.
// Body strings are joined with newlines.
func NewResponse(code int, body ...string) *Response {
	msg, ok := status[code]
	if !ok {
		return NewResponse(500)
	}
	b := fmt.Sprintf("%d %v\r\n", code, msg)
	if len(body) > 0 {
		b = strings.Join(body, "\r\n") + "\r\n"
	}
	return &Response{
		Status:     msg,
		StatusCode: code,
		Headers:    &Headers{},
		Body:       strings.NewReader(b),
	}
}

// Close closes the response's connection.
func (r *Response) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

// Fprint writes the formatted response to w.
func (r *Response) Fprint(w io.Writer) error {
	// Write response status line.
	_, err := fmt.Fprintf(w, "HTTP/1.0 %v %v\r\n", r.StatusCode, r.Status)
	if err != nil {
		return fmt.Errorf("write status line: %v", err)
	}

	// Write header lines.
	if r.Headers != nil {
		err = r.Headers.Fprint(w)
		if err != nil {
			return fmt.Errorf("write headers: %v", err)
		}
	}

	// Write empty line to signal end of headers.
	_, err = fmt.Fprintf(w, "\r\n")
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}

	// Write response body.
	if r.Body != nil {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			return fmt.Errorf("write data: %v", err)
		}
	}

	return nil
}

// ReadResponse parses and reads a response from r.
func ReadResponse(r io.Reader) (*Response, error) {
	res := &Response{
		Headers: &Headers{},
	}

	// Reader is used to read response line by line.
	reader := bufio.NewReader(r)

	// Response body is read from the remaining data in the reader.
	res.Body = reader

	// Read first line of response (status line).
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read status line: %v", err)
	}
	sl := strings.SplitN(strings.TrimSpace(line), " ", 3)
	if len(sl) < 3 {
		return nil, fmt.Errorf("parse status line: \"%v\"", line)
	}
	res.Version = sl[0]
	res.StatusCode, err = strconv.Atoi(sl[1])
	if err != nil {
		return nil, fmt.Errorf("parse status code: %v", err)
	}
	res.Status = sl[2]

	// Read and parse header lines.
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read header line: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		res.Headers.AddRaw(line)
	}

	return res, nil
}
