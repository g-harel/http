package httpc

import (
	"fmt"
	"io"
)

// Headers is a type representing an http request's headers.
type Headers struct {
	data []string
}

// Add adds a name and value combination to the Headers data structure.
func (h *Headers) Add(s string) {
	h.data = append(h.data, s)
}

// Fprint writes the headers to the given writer.
func (h *Headers) Fprint(w io.Writer) error {
	for _, s := range h.data {
		_, err := fmt.Fprintf(w, "%v\r\n", s)
		if err != nil {
			return err
		}
	}
	return nil
}
