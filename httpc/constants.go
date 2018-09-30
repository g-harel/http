package main

const (
	cmdHelp = "help"
	cmdGet  = "get"
	cmdPost = "post"
)

const (
	flagVerbose = "-v"
	flagHeader  = "-h"
	flagData    = "-d"
	flagFile    = "-f"
	flagOut     = "-o"
)

const (
	errMissingCmd  = "no command specified\n"
	errUnknownCmd  = "unknown command\n"
	errTooManyArgs = "too many arguments\n"
	errMissingURL  = "missing url\n"
	errDataAndFile = "cannot use both \"" + flagData + "\" and \"" + flagFile + "\"\n"
	errBadFile     = "could not read file contents"
	errBadOut      = "could not write to output file\n"
)

const (
	msgHelp = `httpc is a curl-like application but supports HTTP protocol only.
Usage:
   httpc command [arguments]
The commands are:
   ` + cmdGet + `     executes a HTTP GET request and prints the response.
   ` + cmdPost + `    executes a HTTP POST request and prints the response.
   ` + cmdHelp + `    prints this screen.
Use "httpc ` + cmdHelp + ` [command]" for more information about a command.`

	msgGet = `usage: httpc ` + cmdGet + ` [` + flagVerbose + `] [` + flagHeader + ` key:value] [` + flagOut + ` filename] URL

Executes a HTTP GET request for a given URL.

   ` + flagVerbose + `             Prints the detail of the response such as protocol, status, and headers.
   ` + flagHeader + ` key:value   Associates headers to HTTP Request with the format 'key:value'.
   ` + flagOut + ` filename    Writes response body to specified file.`

	msgPost = `usage: httpc ` + cmdPost + ` [` + flagVerbose + `] [` + flagHeader +
		` key:value] [` + flagData + ` inline-data] [` + flagFile +
		` file] [` + flagOut + ` filename] URL

Executes a HTTP POST request for a given URL with inline data or from file.

   ` + flagVerbose + `             Prints the detail of the response such as protocol, status, and headers.
   ` + flagHeader + ` key:value   Associates headers to HTTP Request with the format 'key:value'.
   ` + flagData + ` string      Associates an inline data to the body HTTP POST request.
   ` + flagFile + ` file        Associates the content of a file to the body HTTP POST request.
   ` + flagOut + ` filename    Writes response body to specified file.

Either [` + flagData + `] or [` + flagFile + `] can be used but not both.`
)
