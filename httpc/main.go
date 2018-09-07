package main

import (
	"fmt"
	"os"

	"github.com/g-harel/httpc"
)

func main() {
	args := httpc.NewArgs(os.Args[1:])

	h := args.MultiString("-h")

	headers := httpc.Headers{}
	for _, s := range h {
		headers.Add(s)
	}
	fmt.Println(h)
}
