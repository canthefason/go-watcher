package main

import (
	"os"

	watcher "github.com/canthefason/go-watcher"
)

func main() {
	params := watcher.PrepareArgs(os.Args)

	w := watcher.MustRegisterWatcher(params)
	defer w.Close()

	r := watcher.NewRunner()

	// wait for build and run the binary with given params
	go r.Init(params)
	b := watcher.NewBuilder(w, r)

	// build given package
	go b.Build(params)

	// force update for initial package build
	go w.ForceUpdate()

	// listen for further changes
	go w.ListenChanges()

	r.Wait()
}
