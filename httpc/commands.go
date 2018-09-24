package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/httpc"
)

func help(args *Arguments, log *Logger) {
	log.verbose = true
	if args.Match(cmdGet) {
		log.Print(msgGet)
		return
	}
	if args.Match(cmdPost) {
		log.Print(msgPost)
		return
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgHelp)
	}
	log.Print(msgHelp)
}

func get(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)

	req := &httpc.Request{
		Verb:    "GET",
		URL:     "",
		Headers: (&httpc.Headers{}).AddRaw(readHeaders(args)...),
		Data:    nil,
	}

	var ok bool
	req.URL, ok = args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgGet)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgGet)
	}

	err := httpc.HTTP(req, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func post(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)

	req := &httpc.Request{
		Verb:    "POST",
		URL:     "",
		Headers: (&httpc.Headers{}).AddRaw(readHeaders(args)...),
		Data:    nil,
		Len:     0,
		Type:    "",
	}

	d, ok := args.MatchBefore(flagData)
	if ok {
		req.Data = strings.NewReader(d)
		req.Len = len(d)
		req.Type = mime.TypeByExtension(".txt")
	}

	f, ok := args.MatchBefore(flagFile)
	if req.Data != nil && ok {
		log.Fatal(errDataAndFile, msgPost)
	}
	if ok {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(fmt.Sprintf("%v: %v", errBadFile, err))
		}
		s, err := file.Stat()
		if err != nil {
			log.Fatal(fmt.Sprintf("%v: %v", errBadFile, err))
		}
		defer file.Close()

		req.Data = file
		req.Len = int(s.Size())
		req.Type = mime.TypeByExtension(filepath.Ext(f))
	}

	req.URL, ok = args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgPost)
	}

	err := httpc.HTTP(req, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}
