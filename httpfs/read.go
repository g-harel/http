package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/g-harel/http"
)

func read(path string) (*http.Response, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return http.NewResponse(404), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	s, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("get file info: %v", err)
	}
	if s.IsDir() {
		return http.NewResponse(404), nil
	}

	res := http.NewResponse(200)
	res.Headers.Add("Content-Type", mime.TypeByExtension(filepath.Ext(f.Name())))
	res.Body = f

	return res, nil
}
