package main

// Arguments represents a list of arguments being parsed.
type Arguments struct {
	data []string
}

// NewArgs creates a list of args from a given string slice.
func NewArgs(args []string) *Arguments {
	return &Arguments{
		data: args,
	}
}

// Next returns the next argument. If there is none, the returned bool will be false.
func (args *Arguments) Next() (string, bool) {
	if len(args.data) == 0 {
		return "", false
	}
	r := args.data[0]
	args.data = args.data[1:]
	return r, true
}

// Match checks that the first argument is equal to the input strings.
// The matched argument is removed from the slice.
func (args *Arguments) Match(s string) bool {
	if len(args.data) == 0 {
		return false
	}
	r := args.data[0] == s
	if r {
		args.data = args.data[1:]
	}
	return r
}

// MatchBefore checks that the first arg is equal to the input string and returns the next arg.
// Both the matched argument and the one following it are removed from the slice.
// The returned boolean details whether the input string was matched.
func (args *Arguments) MatchBefore(s string) (string, bool) {
	if len(args.data) < 2 {
		return "", false
	}
	if args.data[0] != s {
		return "", false
	}
	r := args.data[1]
	args.data = args.data[2:]
	return r, true
}
