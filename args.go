package httpc

type arg struct {
	value    string
	consumed bool
}

// Args represents a list of arguments being parsed.
type Args []*arg

// NewArgs creates a List from a given string slice.
func NewArgs(s []string) Args {
	l := make(Args, len(s))
	for i, a := range s {
		l[i] = &arg{
			value:    a,
			consumed: false,
		}
	}
	return l
}

// Commands returns a list of all the non-consumed args in the list.
// The values are not processed or filtered beyond their consumed status.
func (l Args) Commands() []string {
	list := []string{}
	for _, a := range l {
		if !a.consumed {
			list = append(list, a.value)
		}
	}
	return list
}

// Bool checks if the given string is in the args.
func (l Args) Bool(f string) bool {
	for _, a := range l {
		if a.value == f {
			a.consumed = true
			return true
		}
	}
	return false
}

// String returns the arg immediately after the first instance of the given string.
func (l Args) String(s string) string {
	for i, a := range l {
		if a.value != s {
			continue
		}
		a.consumed = true
		if i+1 >= len(l) {
			continue
		}
		data := l[i+1]
		data.consumed = true
		return data.value
	}
	return ""
}

// MultiString creates a slice of all the args preceeded by the given string.
func (l Args) MultiString(s string) []string {
	list := []string{}
	for i := 0; i < len(l); i++ {
		match := l[i]
		if match.value != s {
			continue
		}
		match.consumed = true
		if i+1 >= len(l) {
			continue
		}
		i++
		data := l[i]
		data.consumed = true
		list = append(list, data.value)
	}
	return list
}
