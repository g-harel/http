package main

import (
	"fmt"
	"io"
	"os"
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
	h := readHeaders(args)

	url, ok := args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgGet)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgGet)
	}

	req := &httpc.Request{
		Verb:    "GET",
		URL:     url,
		Headers: h,
		Data:    nil,
	}

	err := httpc.HTTP(req, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func post(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)
	h := readHeaders(args)

	var data io.Reader
	d, ok := args.MatchBefore(flagData)
	if ok {
		data = strings.NewReader(d)
	}

	f, ok := args.MatchBefore(flagFile)
	if data != nil && ok {
		log.Fatal(errDataAndFile, msgPost)
	}
	if ok {
		f, err := os.Open(f)
		if err != nil {
			log.Fatal(fmt.Sprintf("%v: %v", errBadFile, err))
		}
		defer f.Close()
		data = f
	}

	url, ok := args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgPost)
	}

	req := &httpc.Request{
		Verb:    "POST",
		URL:     url,
		Headers: h,
		Data:    data,
	}

	err := httpc.HTTP(req, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}
