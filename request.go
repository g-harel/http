package http

import (
	"fmt"
	"io"
	"net/url"
)

// Request represents an HTTP request to be sent.
type Request struct {
	Method  string
	URL     string
	Headers *Headers
	Body    io.Reader

	url *url.URL
}

// Validate checks that the request is valid and computes the request's url.
func (r *Request) Validate() error {
	if r.Method == "GET" && r.Body != nil {
		return fmt.Errorf("Cannot write body to GET request")
	}

	u, err := url.Parse(r.URL)
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
	r.url = u

	return nil
}

// Fprint writes the a formatted request to w.
func (r *Request) Fprint(w io.Writer) error {
	if r.url == nil {
		return fmt.Errorf("cannot print non-validated request")
	}

	// Host header is written with the value extracted from the url.
	r.Headers.Add("Host", r.url.Hostname())

	// Write request status line.
	_, err := fmt.Fprintf(w, "%v %v HTTP/1.0\r\n", r.Method, r.url.RequestURI())
	if err != nil {
		return fmt.Errorf("could not write request line: %v", err)
	}

	// Write header lines.
	err = r.Headers.Fprint(w)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	// Write empty line to signal end of headers.
	_, err = fmt.Fprintf(w, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write: %v", err)
	}

	// Write request body.
	if r.Body != nil {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			return fmt.Errorf("could not write data: %v", err)
		}
	}

	return nil
}
