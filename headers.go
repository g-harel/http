package httpc

import (
	"fmt"
	"io"
	"strings"
)

// Headers is a type representing an http request's headers.
type Headers map[string]string

// Add adds a name and value combination to the Headers data structure.
func (h *Headers) Add(name, value string) {
	(*h)[strings.Title(name)] = value
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
