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
func HTTP(req *Request, log io.Writer) (*Response, error) {
	// Input url is parsed and validated.
	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse given url: %v", err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("missing protocol in \"%v\"", u.String())
	}
	if u.Scheme != "http" {
		return nil, fmt.Errorf("unknown protocol \"%v\" in \"%v\"", u.Scheme, u.String())
	}
	if u.Port() == "" {
		u.Host += ":80"
	}
	if u.Path == "" {
		u.Path = "/"
	}

	// Headers are copied from the given reference.
	h := &Headers{}
	h.AddCopy(req.Headers)

	// Host header is written with the value extracted from the url.
	h.Add("Host", u.Hostname())

	// Open TCP connection to host.
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, fmt.Errorf("could connect to host: %v", err)
	}
	defer conn.Close()

	// All data written to the connection is mirrored into the log writer.
	w := io.MultiWriter(conn, log)

	// Write request status line.
	_, err = fmt.Fprintf(w, "%v %v HTTP/1.0\r\n", req.Method, u.RequestURI())
	if err != nil {
		return nil, fmt.Errorf("could not write request line: %v", err)
	}

	// Write header lines.
	err = h.Fprint(w)
	if err != nil {
		return nil, fmt.Errorf("could not write headers: %v", err)
	}

	// Write empty line to signal end of headers.
	_, err = fmt.Fprintf(w, "\r\n")
	if err != nil {
		return nil, fmt.Errorf("could not write to request: %v", err)
	}

	// Write request body.
	if req.Body != nil {
		_, err := io.Copy(w, req.Body)
		if err != nil {
			return nil, fmt.Errorf("could not write data to request: %v", err)
		}

		// Formatting errors are not critical.
		_, _ = fmt.Fprintf(log, "\n")
	}

	res := &Response{
		Headers: &Headers{},
	}

	// Reader is used to read response line by line.
	reader := bufio.NewReader(conn)
	isEOF := false

	// Response body is read from the remaining data in the reader.
	res.Body = reader

	// Read first line of response (status line).
	line, err := reader.ReadString('\n')
	if err == io.EOF {
		isEOF = true
	} else if err != nil {
		return nil, fmt.Errorf("could not read response status line: %v", err)
	}
	_, err = fmt.Fprintf(log, line)
	if err != nil {
		return nil, fmt.Errorf("could not copy response line to output log: %v", err)
	}
	sl := strings.Split(line, " ")
	if len(sl) < 2 {
		return nil, fmt.Errorf("could not parse response status line: %v", line)
	}
	res.StatusCode, err = strconv.Atoi(sl[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse response status code: %v", err)
	}

	// If the first read has consumed the entire response, there is no need to proceed.
	if isEOF {
		return res, nil
	}

	// Read and parse header lines.
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return res, nil
		} else if err != nil {
			return nil, fmt.Errorf("could not read header line: %v", err)
		}
		_, err = fmt.Fprintf(log, line)
		if err != nil {
			return nil, fmt.Errorf("could not copy response line to output log: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		res.Headers.AddRaw(line)
	}

	return res, nil
}
