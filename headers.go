package httpc

import (
	"strings"
)

// Headers is a type representing an http request's headers.
type Headers map[string]string

// AddString parses the input string to extract and add a name/value combination.
// The format of the string should be "{name}:{value}" where everything before the
// first colon is the name and everything after is the value.
func (h *Headers) AddString(s string) {
	split := strings.SplitN(s, ":", 2)
	name := split[0]
	value := ""
	if len(split) > 1 {
		value = split[1]
	}
	(*h)[name] = value
}

// AddPair adds a name and value combination to the Headers data structure.
func (h *Headers) AddPair(name, value string) {
	(*h)[name] = value
}
