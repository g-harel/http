package httpc

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

// Get makes an http GET request to the provided url.
func Get(addr string, headers []string, log io.Writer, res io.Writer) error {
	u, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("could not parse given url: %v", err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("missing protocol in \"%v\"", u.String())
	}
	if u.Scheme != "http" {
		return fmt.Errorf("unknown protocol \"%v\" in \"%v\"", u.Scheme, u.String())
	}
	headers = append(headers, "Host: "+u.Host)
	if u.Port() == "" {
		u.Host += ":80"
	}
	if u.Path == "" {
		u.Path = "/"
	}

	headers = append(headers, "User-Agent: httpc")

	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return fmt.Errorf("could connect to host: %v", err)
	}
	defer conn.Close()

	req := io.MultiWriter(conn, log)

	_, err = fmt.Fprintf(req, "GET %v HTTP/1.0\r\n", u.RequestURI())
	if err != nil {
		return fmt.Errorf("could not write request line: %v", err)
	}

	for _, s := range headers {
		_, err := fmt.Fprintf(req, "%v\r\n", s)
		if err != nil {
			return fmt.Errorf("could not write headers: %v", err)
		}
	}

	_, err = fmt.Fprintf(req, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write newline: %v", err)
	}

	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err == nil {
			fmt.Fprint(res, string(buf[:n]))
			continue
		}
		if err == io.EOF {
			fmt.Fprint(res, string(buf[:n]))
			break
		}
		return fmt.Errorf("error reading response: %v", err)
	}

	return nil
}

// Post makes an http POST request to the provided url.
func Post(url string, headers []string, data string, log io.Writer, res io.Writer) error {
	fmt.Fprintln(log, "post "+url)
	fmt.Fprintln(log, "data "+data)
	return nil
}
