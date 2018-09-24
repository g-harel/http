package httpc

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

// Request represents an HTTP request to be sent.
type Request struct {
	Verb    string
	URL     string
	Headers *Headers
	Body    io.Reader
}

// Response represents a received HTTP response.
type Response struct {
	StatusCode int
	Headers    *Headers
	Body       io.Reader
}

// HTTP executes an HTTP request.
func HTTP(req *Request, log io.Writer, res io.Writer) error {
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

	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return fmt.Errorf("could connect to host: %v", err)
	}
	defer conn.Close()

	w := io.MultiWriter(conn, log)

	_, err = fmt.Fprintf(w, "%v %v HTTP/1.0\r\n", req.Verb, u.RequestURI())
	if err != nil {
		return fmt.Errorf("could not write request line: %v", err)
	}

	// Headers are copied from input.
	// Host header's value is replaced with the value from the url.
	h := &Headers{}
	h.AddCopy(req.Headers)
	h.Add("Host", u.Hostname())

	err = h.Fprint(w)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	_, err = fmt.Fprintf(w, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write newline: %v", err)
	}

	if req.Body != nil {
		_, err := io.Copy(w, req.Body)
		if err != nil {
			return fmt.Errorf("could not write data: %v", err)
		}

		// Formatting errors are not critical.
		fmt.Fprintf(log, "\n\n")
	}

	_, err = io.Copy(res, conn)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Formatting errors are not critical.
	fmt.Fprintf(log, "\n")

	return nil
}
