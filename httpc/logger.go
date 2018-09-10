package main

import (
	"fmt"
	"os"
)

// Logger is a helper which handles all logging logic.
type Logger struct {
	verbose bool
}

// Message only prints the input string if the logger is verbose.
func (l *Logger) Message(str string) {
	if l.verbose {
		fmt.Println(str)
	}
}

// Write implements "io.Writer" and respects the verbose option.
func (l *Logger) Write(b []byte) (int, error) {
	if l.verbose {
		return fmt.Fprintf(os.Stdin, string(b))
	}
	return len(b), nil
}

// Result prints the input string and exits with no error.
func (l *Logger) Result(str string) {
	fmt.Println(str)
	os.Exit(0)
}

// Fatal prints the input string and exits with an error code.
func (l *Logger) Fatal(str string) {
	fmt.Fprintln(os.Stderr, str)
	os.Exit(1)
}
