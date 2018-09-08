package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/g-harel/httpc"
)

var helpMsg = `httpc is a curl-like application but supports HTTP protocol only.
Usage:
   httpc command [arguments]
The commands are:
   get     executes a HTTP GET request and prints the response.
   post    executes a HTTP POST request and prints the response.
   help    prints this screen.
Use "httpc help [command]" for more information about a command.`

var helpGetMsg = `usage: httpc get [-v] [-h key:value] URL

Get executes a HTTP GET request for a given URL.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.`

var helpPostMsg = `usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL

Post executes a HTTP POST request for a given URL with inline data or from file.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.
   -d string      Associates an inline data to the body HTTP POST request.
   -f file        Associates the content of a file to the body HTTP POST request.

Either [-d] or [-f] can be used but not both.`

func main() {
	args := NewArgs(os.Args[1:])
	log := Logger{}

	if args.Match([]string{"help", "get"}) {
		log.Result(helpGetMsg)
	}

	if args.Match([]string{"help", "post"}) {
		log.Result(helpPostMsg)
	}

	if args.Match([]string{"help"}) {
		log.Result(helpMsg)
	}

	log.verbose = args.Bool("-v")

	headers := httpc.Headers{}
	for _, s := range args.MultiString("-h") {
		split := strings.SplitN(s, ":", 2)
		name := split[0]
		value := ""
		if len(split) > 1 {
			value = split[1]
		}
		headers.AddPair(name, value)
	}

	if args.Match([]string{"get"}) {
		u := args.Unused()
		if len(u) != 1 {
			log.Fatal(helpGetMsg)
		}
		res, err := httpc.Get(u[0], &headers, log.Message)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Result(res)
	}

	if args.Match([]string{"post"}) {
		if len(args.Unused()) != 1 {
			log.Fatal(helpPostMsg)
		}
		fmt.Printf("post => %v\n", args.Unused())
		return
	}

	log.Fatal(helpMsg)
}
