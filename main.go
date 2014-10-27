package main

import "os"

const binaryName = "goldorf-main"

func main() {
	w := MustRegisterWatcher()
	defer w.Close()

	done := make(chan struct{})

	r := NewRunner()

	// first build given package
	go build(w, r)

	// run the binary with given arguments
	go r.Init(getArgs()...)

	// listen for further changes
	go w.ListenChanges()

	<-done
}

func getArgs() []string {
	args := os.Args

	return args[1:len(args)]
}
