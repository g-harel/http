package main

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/g-harel/http"
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

	req := &http.Request{
		Method:  "GET",
		URL:     "",
		Headers: (&http.Headers{}).AddRaw(readHeaders(args)...),
		Body:    nil,
	}

	var file *os.File
	var err error
	filename, ok := args.MatchBefore(flagOut)
	if ok {
		file, err = os.Create(filename)
		if err != nil {
			log.Fatal(errBadOut, msgGet)
		}
	}

	req.URL, ok = args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgGet)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgGet)
	}

	client := http.Client{
		FollowRedirect: true,
		Logger:         log,
	}
	res, err := client.Send(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Close()

	if file != nil {
		_, err := io.Copy(file, res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Print(fmt.Sprintf("Output written to %v", filename))
	} else {
		_, err := io.Copy(os.Stdout, res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func post(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)

	headers := &http.Headers{}
	headers.AddRaw(readHeaders(args)...)

	req := &http.Request{
		Method:  "POST",
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

	var file *os.File
	var err error
	filename, ok := args.MatchBefore(flagOut)
	if ok {
		file, err = os.Create(filename)
		if err != nil {
			log.Fatal(errBadOut, msgGet)
		}
	}

	req.URL, ok = args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgPost)
	}

	client := http.Client{
		FollowRedirect: true,
		Logger:         log,
	}
	res, err := client.Send(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Close()

	if file != nil {
		_, err := io.Copy(file, res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Print(fmt.Sprintf("Output written to %v", filename))
	} else {
		_, err := io.Copy(os.Stdout, res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
