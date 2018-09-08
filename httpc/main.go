package main

import (
	"fmt"
	"os"

	"github.com/g-harel/httpc"
)

var helpMsg = `
httpc is a curl-like application but supports HTTP protocol only.
Usage:
   httpc command [arguments]
The commands are:
   get     executes a HTTP GET request and prints the response.
   post    executes a HTTP POST request and prints the response.
   help    prints this screen.
Use "httpc help [command]" for more information about a command.`

var helpGetMsg = `
usage: httpc get [-v] [-h key:value] URL

Get executes a HTTP GET request for a given URL.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.`

var helpPostMsg = `
usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL

Post executes a HTTP POST request for a given URL with inline data or from file.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.
   -d string      Associates an inline data to the body HTTP POST request.
   -f file        Associates the content of a file to the body HTTP POST request.

Either [-d] or [-f] can be used but not both.`

func handleHelpMsg(msg string, clean bool) {
	fmt.Println(msg + "\n")
	if clean {
		os.Exit(0)
	}
	os.Exit(1)
}

func main() {
	args := NewArgs(os.Args[1:])

	if args.Match([]string{"help", "get"}) {
		handleHelpMsg(helpGetMsg, true)
	}

	if args.Match([]string{"help", "post"}) {
		handleHelpMsg(helpPostMsg, true)
	}

	if args.Match([]string{"help"}) {
		handleHelpMsg(helpMsg, true)
	}

	verbose := args.Bool("-v")

	headers := httpc.Headers{}
	for _, s := range args.MultiString("-h") {
		headers.AddString(s)
	}

	if args.Match([]string{"get"}) {
		if len(args.Unused()) != 1 {
			handleHelpMsg(helpGetMsg, false)
		}
		fmt.Printf("get => %v\n", args.Unused())
		return
	}

	if args.Match([]string{"post"}) {
		if len(args.Unused()) != 1 {
			handleHelpMsg(helpPostMsg, false)
		}
		fmt.Printf("post => %v\n", args.Unused())
		return
	}

	fmt.Println(verbose, headers, args.Unused())

	handleHelpMsg(helpMsg, false)
}
