package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sioncojp/flac2float"
)

const (
	app = "flac2float"
)

type CommandOpts struct {
	Interval uint `long:"interval" short:"i" default:"1" description:"It can output specified numerical multiple times finely"`
}

func main() {
	opts := CommandOpts{}
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	log.SetOutput(os.Stderr)
	log.SetPrefix(app + ": ")

	decode := flac2float.New(os.Stdin, opts.Interval)
	values, err := decode.ReadSound()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	for _, value := range values {
		fmt.Println(value)
	}
}
