package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/g-harel/http"
)

func write(path string, r io.Reader, size int64) (*http.Response, error) {
	f, err := os.Create(path)
	if err != nil {
		io.CopyN(ioutil.Discard, r, size)
		if err.(*os.PathError).Err == syscall.EISDIR {
			return http.NewResponse(400, "Cannot Write To Directory"), nil
		}
		return nil, fmt.Errorf("create file: %v", err)
	}

	_, err = io.CopyN(f, r, size)
	if err != nil {
		return nil, fmt.Errorf("write to file: %v", err)
	}

	return http.NewResponse(201), nil
}
