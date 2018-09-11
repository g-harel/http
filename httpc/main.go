package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/g-harel/httpc"
)

var cmdHelp = "help"
var cmdGet = "get"
var cmdPost = "post"

var flagVerbose = "-v"
var flagHeader = "-h"
var flagData = "-d"
var flagFile = "-f"

var errMissingCmd = "No Command Specified\n"
var errUnknownCmd = "Unknown Command\n"
var errTooManyArgs = "Too Many Arguments\n"
var errMissingURL = "Missing URL\n"
var errDataAndFile = "Cannot Use Both\"-d\" and \"-f\"\n"
var errBadFile = "Could Not Read File Contents"

var msgHelp = `httpc is a curl-like application but supports HTTP protocol only.
Usage:
   httpc command [arguments]
The commands are:
   get     executes a HTTP GET request and prints the response.
   post    executes a HTTP POST request and prints the response.
   help    prints this screen.
Use "httpc help [command]" for more information about a command.`

var msgGet = `usage: httpc get [-v] [-h key:value] URL

Get executes a HTTP GET request for a given URL.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.`

var msgPost = `usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL

Post executes a HTTP POST request for a given URL with inline data or from file.

   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.
   -d string      Associates an inline data to the body HTTP POST request.
   -f file        Associates the content of a file to the body HTTP POST request.

Either [-d] or [-f] can be used but not both.`

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

func help(args *Arguments, log *Logger) {
	log.verbose = true
	if args.Match(cmdGet) {
		log.Print(msgGet)
	}
	if args.Match(cmdPost) {
		log.Print(msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgHelp)
	}
	log.Print(msgHelp)
}

func get(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)
	h := readHeaders(args)

	url, ok := args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgGet)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgGet)
	}

	err := httpc.Get(url, h, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func post(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)
	h := readHeaders(args)

	data := ""
	d, dataOK := args.MatchBefore(flagData)
	if dataOK {
		data = d
	}

	f, fileOK := args.MatchBefore(flagFile)
	if fileOK && dataOK {
		log.Fatal(errDataAndFile, msgPost)
	}
	if fileOK {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(fmt.Sprintf("%v: %v", errBadFile, err))
		}
		data = string(file)
	}

	url, ok := args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgPost)
	}

	err := httpc.Post(url, h, data, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
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
