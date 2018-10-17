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

// Close closes the response's body.
func (r *Response) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
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
		return nil, fmt.Errorf("could not read status line: %v", err)
	}
	sl := strings.Split(strings.TrimSpace(line), " ")
	if len(sl) < 3 {
		return nil, fmt.Errorf("could not parse status line: %v", line)
	}
	res.Version = sl[0]
	res.StatusCode, err = strconv.Atoi(sl[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse status code: %v", err)
	}
	res.Status = sl[2]

	// Read and parse header lines.
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read header line: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		res.Headers.AddRaw(line)
	}

	return res, nil
}
