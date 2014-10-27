package main

import (
	"log"
	"os"
)

const binaryName = "goldorf-main"

var (
	done chan struct{}
)

func initChannels() {
	done = make(chan struct{})
}

func main() {
	w := MustRegisterWatcher()
	defer w.Close()
	defer close(done)

	r := NewRunner()

	initChannels()

	// first build given package
	go build(w, r)

	// run the binary with given arguments
	go r.Init(getArgs()...)

	// listen for further changes
	go w.ListenChanges()

	<-done

	log.Print("exiting")
}

func getArgs() []string {
	args := os.Args

	return args[1:len(args)]
}
