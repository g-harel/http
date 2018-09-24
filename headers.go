package httpc

import (
	"fmt"
	"io"
	"strings"
)

// Headers is a type representing an http request's headers.
type Headers map[string]string

// Add adds a name/value combination to the Headers data structure.
func (h *Headers) Add(name, value string) {
	name = strings.TrimSpace(name)
	(*h)[name] = value
}

// AddRaw parses the input string to extract and add name/value combinations.
// The format of the string should be "{name}:{value}".
// Everything before the first colon is the name and everything after is the value.
func (h *Headers) AddRaw(s string) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		split := strings.SplitN(line, ":", 2)

		key := strings.TrimSpace(split[0])
		if key == "" {
			continue
		}

		value := ""
		if len(split) > 1 {
			value = split[1]
		}

		(*h)[key] = value
	}
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
