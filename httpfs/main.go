package main

import (
	"fmt"

	"github.com/g-harel/http"
)

func main() {
	server := http.Server{
		Errors: make(chan error),
	}

	server.Use(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("test")
	})

	go func() {
		for {
			err := <-server.Errors
			println(err.Error())
		}
	}()

	err := server.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
