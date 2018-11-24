package http

import (
	"fmt"
	"io"

	"github.com/g-harel/http/transport"
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
	conn, err := transport.Dial(transportProtocol, fmt.Sprintf("%v:%v", req.Hostname, port))
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

	// Response is written to the log.
	if req.Body != nil {
		_, _ = fmt.Fprintf(c.Logger, "\n")
	}
	_, _ = fmt.Fprintf(c.Logger, "%v %v %v\n", res.Version, res.StatusCode, res.Status)
	_ = res.Headers.Fprint(c.Logger)
	_, _ = fmt.Fprintf(c.Logger, "\n")

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
