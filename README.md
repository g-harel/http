<!--

TODO
- add timeouts for both client and server implementations
- add worker pool to handle requests instead of being sequential

-->

# http

[![](https://godoc.org/github.com/g-harel/http?status.svg)](http://godoc.org/github.com/g-harel/http)

_This project contains solutions for the [first two assignments](./assignments) for the `Data Communications & Computer Networks` course._

The parent directory is a package exposing partial `http/1.0` client and server implementation similar to [`net/http`](https://golang.org/pkg/net/http/).

The [`./httpc`](./httpc) and [`./httpfs`](./https) directories contain a _curl-like_ command line tool and a simple file server respectively.

## httpc

```
$ go get -u github.com/g-harel/http/httpc
```

#### get

```
Executes a HTTP GET request for a given URL.

Usage:
    httpc get [-v] [-h key:value] [-o filename] URL

Flags:
   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.
   -o filename    Writes response body to specified file.
```

#### post

```
Executes a HTTP POST request for a given URL with inline data or from file.

Usage:
   httpc post [-v] [-h key:value] [-d inline-data] [-f file] [-o filename] URL

Flags:
   -v             Prints the detail of the response such as protocol, status, and headers.
   -h key:value   Associates headers to HTTP Request with the format 'key:value'.
   -d string      Associates an inline data to the body HTTP POST request.
   -f filename    Associates the content of a file to the body HTTP POST request.
   -o filename    Writes response body to specified file.
```

_Either `-d` or `-f` can be used but not both._

## httpfs

```
$ go get -u github.com/g-harel/http/httpfs
```

```
Starts a simple file server.

Usage:
   httpfs [-v] [-p port] [-d directory]

Flags:
   -v             Prints debugging messages.
   -p port        Specifies the port for the server to listen to (default "8080").
   -d directory   Specifies the file server's root directory (default ".").
```
