package main

import (
	"fmt"
	"io/ioutil"
	"os"

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

	err := httpc.Get(url, h, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func post(args *Arguments, log *Logger) {
	log.verbose = args.Match(flagVerbose)
	h := readHeaders(args)

	data := ""
	d, dataOK := args.MatchBefore(flagData)
	if dataOK {
		data = d
	}

	f, fileOK := args.MatchBefore(flagFile)
	if fileOK && dataOK {
		log.Fatal(errDataAndFile, msgPost)
	}
	if fileOK {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(fmt.Sprintf("%v: %v", errBadFile, err))
		}
		data = string(file)
	}

	url, ok := args.Next()
	if !ok {
		log.Fatal(errMissingURL, msgPost)
	}
	if _, ok := args.Next(); ok {
		log.Fatal(errTooManyArgs, msgPost)
	}

	err := httpc.Post(url, h, data, log, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}
