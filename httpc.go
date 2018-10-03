package httpc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Request represents an HTTP request to be sent.
type Request struct {
	Method  string
	URL     string
	url     *url.URL
	Headers *Headers
	Body    io.Reader
}

// Response represents a received HTTP response.
type Response struct {
	Version    string
	Status     string
	StatusCode int
	Headers    *Headers
	Body       io.Reader
	conn       net.Conn
}

// Close closes the connection to the host.
func (r *Response) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

// Validate checks that the input request is valid.
// If the function does not error, the input request's `url` field will have been assigned a value.
func validate(req *Request) error {
	if req.Method == "GET" && req.Body != nil {
		return fmt.Errorf("Cannot write body to GET request")
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("could not parse given url: %v", err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("missing protocol in \"%v\"", u.String())
	}
	if u.Scheme != "http" {
		return fmt.Errorf("unknown protocol \"%v\" in \"%v\"", u.Scheme, u.String())
	}
	if u.Port() == "" {
		u.Host += ":80"
	}
	if u.Path == "" {
		u.Path = "/"
	}

	// Validated URL is added to the request's private field.
	req.url = u

	return nil
}

// Request writes request data to the given connection.
func request(conn io.Writer, req *Request) error {
	// Host header is written with the value extracted from the url.
	req.Headers.Add("Host", req.url.Hostname())

	// Write request status line.
	_, err := fmt.Fprintf(conn, "%v %v HTTP/1.0\r\n", req.Method, req.url.RequestURI())
	if err != nil {
		return fmt.Errorf("could not write request line: %v", err)
	}

	// Write header lines.
	err = req.Headers.Fprint(conn)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	// Write empty line to signal end of headers.
	_, err = fmt.Fprintf(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write: %v", err)
	}

	// Write request body.
	if req.Body != nil {
		_, err := io.Copy(conn, req.Body)
		if err != nil {
			return fmt.Errorf("could not write data: %v", err)
		}
	}

	return nil
}

// Response reads response data from a given connection.
func response(conn io.Reader) (*Response, error) {
	res := &Response{
		Headers: &Headers{},
	}

	// Reader is used to read response line by line.
	reader := bufio.NewReader(conn)

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

// HTTP executes an HTTP request.
func HTTP(req *Request, log io.Writer) (*Response, error) {
	// Validate request.
	err := validate(req)
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
	err = request(w, req)
	if err != nil {
		return nil, fmt.Errorf("could not write request: %v", err)
	}

	// Read and parse response from connection.
	res, err := response(conn)
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
	if res.StatusCode == 301 || res.StatusCode == 302 {
		// Body contents are ignored and connection is closed.
		res.Close()

		location, ok := res.Headers.Read("Location")
		if !ok {
			return nil, fmt.Errorf("could not read redirect location")
		}

		req.URL = location
		return HTTP(req, log)
	}

	return res, nil
}
