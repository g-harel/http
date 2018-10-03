package httpc

import (
	"fmt"
	"io"
	"strings"
)

// Helper to format header names consistently for reads and writes.
func formatName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Title(name)
	return name
}

// Headers type represents an http request's headers.
type Headers map[string]string

// Add adds a name/value combination to the Headers' data.
func (h *Headers) Add(name, value string) *Headers {
	name = formatName(name)
	if name != "" {
		(*h)[name] = value
	}
	return h
}

// AddRaw parses the input string to extract and add name/value combinations.
// The format of the string should be "{name}:{value}".
// Everything before the first colon is the name and everything after is the value.
func (h *Headers) AddRaw(lines ...string) *Headers {
	for _, line := range lines {
		split := strings.SplitN(line, ":", 2)

		value := ""
		if len(split) > 1 {
			value = strings.TrimSpace(split[1])
		}

		h.Add(split[0], value)
	}
	return h
}

// Read reads the header value for the given name.
func (h *Headers) Read(name string) (string, bool) {
	value, ok := (*h)[formatName(name)]
	return value, ok
}

// Fprint writes the headers to the given writer.
func (h *Headers) Fprint(w io.Writer) error {
	for name, value := range *h {
		_, err := fmt.Fprintf(w, "%v: %v\r\n", name, value)
		if err != nil {
			return err
		}
	}
	return nil
}
