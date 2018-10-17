package http

import (
	"fmt"
	"io"
	"net"
)

// Client is used to configure and send http requests.
type Client struct {
	// FollowRedirect specifies whether 301 or 302 response are followed or not.
	FollowRedirect bool
}

// Send sends an http request.
func (c *Client) Send(req *Request, log io.Writer) (*Response, error) {
	// Validate request.
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("could validate request: %v", err)
	}

	// Open TCP connection to host.
	conn, err := net.Dial("tcp", req.url.Host)
	if err != nil {
		return nil, fmt.Errorf("could connect to host: %v", err)
	}

	// All data written to the connection is mirrored into the log.
	w := io.MultiWriter(conn, log)

	// Write request to connection.
	err = req.Fprint(w)
	if err != nil {
		return nil, fmt.Errorf("could not write request: %v", err)
	}

	// Read and parse response from connection.
	res, err := ReadResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	res.conn = conn

	// Response status line and headers are written to the log.
	if req.Body != nil {
		_, err = fmt.Fprintf(log, "\n")
		if err != nil {
			return nil, fmt.Errorf("could write to log: %v", err)
		}
	}
	_, err = fmt.Fprintf(log, "%v %v %v\n", res.Version, res.StatusCode, res.Status)
	if err != nil {
		return nil, fmt.Errorf("could not log response status line: %v", err)
	}
	err = res.Headers.Fprint(log)
	if err != nil {
		return nil, fmt.Errorf("could not log response headers: %v", err)
	}
	_, err = fmt.Fprintf(log, "\n")
	if err != nil {
		return nil, fmt.Errorf("could write to log: %v", err)
	}

	// Redirects are followed without intervention.
	if c.FollowRedirect && (res.StatusCode == 301 || res.StatusCode == 302) {
		// Body contents are ignored and connection is closed.
		res.Close()

		location, ok := res.Headers.Read("Location")
		if !ok {
			return nil, fmt.Errorf("could not read redirect location")
		}

		req.URL = location
		return c.Send(req, log)
	}

	return res, nil
}