package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
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
		Body:    nil,
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

	headers := &httpc.Headers{}
	headers.AddRaw(readHeaders(args)...)

	req := &httpc.Request{
		Verb:    "POST",
		URL:     "",
		Headers: headers,
		Body:    nil,
	}

	d, ok := args.MatchBefore(flagData)
	if ok {
		req.Body = strings.NewReader(d)
		headers.Add("Content-Length", strconv.Itoa(len(d)))
		headers.Add("Content-Type", mime.TypeByExtension(".txt"))
	}

	f, ok := args.MatchBefore(flagFile)
	if req.Body != nil && ok {
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

		req.Body = file
		headers.Add("Content-Length", strconv.Itoa(int(s.Size())))
		headers.Add("Content-Type", mime.TypeByExtension(filepath.Ext(f)))
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
