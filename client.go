package http

import (
	"fmt"
	"io"
	"net"
)

// Client is used to configure and send http requests.
type Client struct {
	FollowRedirect bool
	Logger         io.Writer
}

// Send sends an http request.
func (c *Client) Send(req *Request) (*Response, error) {
	if req.Method == "GET" && req.Body != nil {
		return nil, fmt.Errorf("Cannot write body to GET request")
	}

	// Open TCP connection to host.
	port := req.Port
	if port == "" {
		port = "80"
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", req.Hostname, port))
	if err != nil {
		return nil, fmt.Errorf("could connect to host: %v", err)
	}

	// All data written to the connection is mirrored into the log.
	w := io.MultiWriter(conn, c.Logger)

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
		_, err = fmt.Fprintf(c.Logger, "\n")
		if err != nil {
			return nil, fmt.Errorf("could write to log: %v", err)
		}
	}
	_, err = fmt.Fprintf(c.Logger, "%v %v %v\n", res.Version, res.StatusCode, res.Status)
	if err != nil {
		return nil, fmt.Errorf("could not log response status line: %v", err)
	}
	err = res.Headers.Fprint(c.Logger)
	if err != nil {
		return nil, fmt.Errorf("could not log response headers: %v", err)
	}
	_, err = fmt.Fprintf(c.Logger, "\n")
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

		err = req.URL(location)
		if err != nil {
			return nil, fmt.Errorf("could not parse redirect url: %v", err)
		}

		return c.Send(req)
	}

	return res, nil
}
