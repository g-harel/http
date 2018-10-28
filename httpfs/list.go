package main

import (
	"fmt"
	"os"

	"github.com/g-harel/http"
)

func list(dir string) (*http.Response, error) {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return http.NewResponse(404), nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	s, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get file info: %v", err)
	}
	if !s.IsDir() {
		return http.NewResponse(404), nil
	}

	files, err := f.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %v", err)
	}

	filenames := []string{}
	for _, f := range files {
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}
		filenames = append(filenames, name)
	}

	return http.NewResponse(200, filenames...), nil
}
