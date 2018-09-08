package httpc

// Headers is a type representing an http request's headers.
type Headers map[string]string

// AddPair adds a name and value combination to the Headers data structure.
func (h *Headers) AddPair(name, value string) {
	(*h)[name] = value
}
