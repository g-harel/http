package main

import (
	"fmt"
	"os"
)

// Logger is a helper which handles all logging logic.
type Logger struct {
	verbose bool
}

// Error prints the error in red to stderr.
func (l *Logger) Error(err error) {
	if !l.verbose {
		fmt.Fprintf(os.Stderr, "\n\033[31;1m[HTTPFS] %v\033[0m\n", err)
	}
}
