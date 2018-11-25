package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/g-harel/http"
)

var verbose = flag.Bool("v", false, "Prints debugging messages.")
var port = flag.String("p", "8080", "Specifies the port number that the server will listen and serve at. Default is 8080.")
var dir = flag.String("d", "", "Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "httpfs is a simple file server.\nusage: httpfs [-v] [-p PORT] [-d PATH-TO-DIR]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	logger := &Logger{
		verbose: *verbose,
	}

	// Find base path whose contents will be served.
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	base := path.Join(wd, *dir)

	// Create server and configure error logging.
	server := http.Server{
		ErrChan: make(chan error),
	}
	go func() {
		for {
			logger.Error(<-server.ErrChan)
		}
	}()

	// Assign a request handler to the server.
	server.Use(func(req *http.Request) (*http.Response, error) {
		// Check for path-related errors.
		if req.Path == "" {
			return http.NewResponse(400, "Empty Path"), nil
		}
		if req.Path[0] != '/' {
			return http.NewResponse(400, "Malformed Path"), nil
		}

		// List requests end with `/`.
		isListRequest := req.Path[len(req.Path)-1] == '/'

		// Check that the path does not point outside the working directory.
		absolutePath := path.Join(base, req.Path)
		deltaPath, err := filepath.Rel(base, absolutePath)
		if err != nil {
			return nil, fmt.Errorf("check for dangerous path")
		}
		if strings.Index(deltaPath, "../") >= 0 {
			return http.NewResponse(403), nil
		}

		// Handle a file write.
		if req.Method == "POST" {
			s, ok := req.Headers.Read("Content-Length")
			if !ok {
				return http.NewResponse(400, "Missing Content-Length"), nil
			}

			size, err := strconv.Atoi(s)
			if err != nil {
				return http.NewResponse(400, "Invalid Content-Length"), nil
			}

			return write(absolutePath, req.Body, int64(size))
		}

		// Handle a file or directory read.
		if req.Method == "GET" {
			if isListRequest {
				return list(absolutePath)
			}
			return read(absolutePath)
		}

		return http.NewResponse(400), nil
	})

	// Start listening on specified port.
	err = server.Listen(":" + *port)
	if err != nil {
		logger.Fatal(err)
	}
}
