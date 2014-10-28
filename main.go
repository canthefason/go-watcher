package main

import "os"

const binaryName = "goldorf-main"

func main() {
	args, rootpackage := prepareArgs()

	w := MustRegisterWatcher(rootpackage)
	defer w.Close()

	done := make(chan struct{})

	r := NewRunner()

	// first build given package
	go build(w, r)

	// run the binary with given arguments
	go r.Init(args...)
	// force update for initial package build
	go w.ForceUpdate()

	// listen for further changes
	go w.ListenChanges()

	<-done
}

func prepareArgs() ([]string, string) {
	args := os.Args

	// remove command
	args = args[1:len(args)]

	filteredArgs := make([]string, 0)
	var rootpackage string
	// !!! TODO change this part, it is error prone for parameters w/o values.
	for i := 0; i < len(args); i += 2 {
		if args[i] == "-rootpackage" || args[i] == "--rootpackage" {
			rootpackage = args[i+1]
			continue
		}

		filteredArgs = append(filteredArgs, args[i], args[i+1])
	}

	return filteredArgs, rootpackage
}
