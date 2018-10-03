package main

import (
	"fmt"
	"os"
)

// Logger is a helper which handles all logging logic.
type Logger struct {
	verbose bool
}

// Write implements "io.Writer" while respecting the verbose option.
// Output is colored grey.
func (l *Logger) Write(b []byte) (int, error) {
	if l.verbose {
		n, err := fmt.Printf("\033[2m%v\033[0m", string(b))
		return n - 8, err
	}
	return len(b), nil
}

// Print prints the input string.
func (l *Logger) Print(str string) {
	fmt.Println(str)
}

// Fatal prints the input string and exits with an error code.
// Output is colored red.
func (l *Logger) Fatal(strs ...string) {
	for _, str := range strs {
		fmt.Printf("\n\033[31;1m%v\033[0m\n\n", str)
	}
	os.Exit(1)
}
