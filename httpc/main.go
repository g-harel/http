package main

import (
	"os"
)

func readHeaders(args *Arguments) []string {
	headers := []string{}
	for {
		h, ok := args.MatchBefore(flagHeader)
		if !ok {
			break
		}
		headers = append(headers, h)
	}
	return headers
}

func main() {
	args := NewArgs(os.Args[1:])
	log := &Logger{}

	cmd, ok := args.Next()
	if !ok {
		log.Fatal(errMissingCmd, msgHelp)
	}

	switch cmd {
	case cmdHelp:
		help(args, log)
	case cmdGet:
		get(args, log)
	case cmdPost:
		post(args, log)
	default:
		log.Fatal(errUnknownCmd, msgHelp)
	}
}
