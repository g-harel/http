package main

import (
	"github.com/g-harel/http"
)

func main() {
	server := http.Server{}

	err := server.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
