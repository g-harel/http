package httpc

import (
	"strings"
)

// Headers is a type representing an http request's headers.
type Headers map[string]string

// Add parses the input string to extract and add a name/value combination.
// The format of the string should be "{name}:{value}".
// Everything before the first colon is the name and everything after is the value.
func (h *Headers) Add(s string) {
	split := strings.SplitN(s, ":", 2)
	key := split[0]
	value := ""
	if len(split) > 1 {
		value = split[1]
	}
	(*h)[key] = value
}
