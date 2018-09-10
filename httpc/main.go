package main

import (
	"fmt"
	"io/ioutil"
	"os"

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

	if args.Match("help") {
		log.verbose = true
		if args.Match("get") {
			log.Result(helpGetMsg)
		}
		if args.Match("post") {
			log.Result(helpPostMsg)
		}
		if _, ok := args.Next(); ok {
			log.Message("Error: too many arguments")
			log.Fatal(helpMsg)
		}
		log.Result(helpMsg)
	}

	if args.Match("get") {
		log.verbose = args.Match("-v")
		headers := httpc.Headers{}
		for {
			h, ok := args.MatchBefore("-h")
			if !ok {
				break
			}
			headers.Add(h)
		}

		url, ok := args.Next()
		if !ok {
			log.Message("Error: missing url")
			log.Fatal(helpGetMsg)
		}
		if _, ok := args.Next(); ok {
			log.Message("Error: too many arguments")
			log.Fatal(helpGetMsg)
		}

		err := httpc.Get(url, &headers, &log, os.Stdout)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	if args.Match("post") {
		log.verbose = args.Match("-v")
		headers := httpc.Headers{}
		for {
			h, ok := args.MatchBefore("-h")
			if !ok {
				break
			}
			headers.Add(h)
		}

		data := ""
		d, dataOK := args.MatchBefore("-d")
		if dataOK {
			data = d
		}

		f, fileOK := args.MatchBefore("-f")
		if fileOK && dataOK {
			log.Message("Error: cannot use both '-d' and '-f' flags")
			log.Fatal(helpPostMsg)
		}
		if fileOK {
			file, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal(fmt.Sprintf("Error: could not read file contents: %v", err))
			}
			data = string(file)
		}

		url, ok := args.Next()
		if !ok {
			log.Message("Error: missing url")
			log.Fatal(helpPostMsg)
		}
		if _, ok := args.Next(); ok {
			log.Message("Error: too many arguments")
			log.Fatal(helpPostMsg)
		}

		res, err := httpc.Post(url, &headers, data, log.Message)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Result(res)
		return
	}

	log.Fatal(helpMsg)
}
