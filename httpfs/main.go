package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/g-harel/http"
)

var verbose = flag.Bool("v", false, "Prints debugging messages.")
var port = flag.String("p", "8080", "Specifies the port number that the server will listen and serve at. Default is 8080.")
var dir = flag.String("d", "", "Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.")

func main() {
	flag.Parse()

	// Find base path whose contents will be served.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	base := path.Join(wd, *dir)

	server := http.Server{
		ErrChan: make(chan error),
	}

	go func() {
		for {
			err := <-server.ErrChan
			fmt.Printf("\033[31;1m%v\033[0m\n", err.Error())
		}
	}()

	server.Use(func(req *http.Request) (*http.Response, error) {
		if req.Path == "" {
			return http.NewResponse(400, "Empty Path"), nil
		}
		if req.Path[0] != '/' {
			return http.NewResponse(400, "Malformed Path"), nil
		}
		absolutePath := path.Join(base, path.Clean(req.Path))
		isList := req.Path[len(req.Path)-1] == '/'

		if req.Method == "POST" {
			return http.NewResponse(200), nil
		}

		if req.Method == "GET" {
			if isList {
				return list(absolutePath)
			}
			return http.NewResponse(200), nil
		}

		return http.NewResponse(400), nil
	})

	err = server.Listen(":" + *port)
	if err != nil {
		panic(err)
	}
}
