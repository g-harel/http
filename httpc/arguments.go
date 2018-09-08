package main

// Arguments represents a list of arguments being parsed.
type Arguments []*struct {
	value string
	used  bool
}

// NewArgs creates a list of args from a given string slice.
func NewArgs(args []string) Arguments {
	l := make(Arguments, len(args))
	for i, arg := range args {
		l[i] = &struct {
			value string
			used  bool
		}{
			value: arg,
			used:  false,
		}
	}
	return l
}

// Unused returns a list of all the unused args in the list.
// The values are not processed or filtered beyond their used status.
func (args Arguments) Unused() []string {
	list := []string{}
	for _, arg := range args {
		if !arg.used {
			list = append(list, arg.value)
		}
	}
	return list
}

// Match checks that the first arguments are equal to the input strings.
// If an only if everything matches, the involved arguments are marked as used.
func (args Arguments) Match(strs []string) bool {
	if len(strs) > len(args) {
		return false
	}
	for i, str := range strs {
		if str != args[i].value {
			return false
		}
	}
	for i := 0; i < len(strs); i++ {
		args[i].used = true
	}
	return true
}

// Bool checks if the given string is in the args.
func (args Arguments) Bool(str string) bool {
	for _, arg := range args {
		if arg.value == str {
			arg.used = true
			return true
		}
	}
	return false
}

// String returns the arg immediately after the first instance of the given string.
func (args Arguments) String(str string) string {
	for i, arg := range args {
		if arg.value != str {
			continue
		}
		arg.used = true
		if i+1 >= len(args) {
			continue
		}
		data := args[i+1]
		data.used = true
		return data.value
	}
	return ""
}

// MultiString creates a slice of all the args preceeded by the given string.
func (args Arguments) MultiString(str string) []string {
	list := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg.value != str {
			continue
		}
		arg.used = true
		if i+1 >= len(args) {
			continue
		}
		i++
		data := args[i]
		data.used = true
		list = append(list, data.value)
	}
	return list
}
