// Goldorf is a command line tool inspired by fresh (https://github.com/pilu/fresh) and used
// for watching .go file changes, and restarting the app in case of an update/delete/add operation.
// After you installed it, you can run your apps with their default parameters as:
// goldorf -c config -p 7000 -h localhost
package main

// Binary name used for built package
const binaryName = "goldorf-main"

func main() {
	params := prepareArgs()

	w := MustRegisterWatcher(params)
	defer w.Close()

	done := make(chan struct{})

	r := NewRunner()

	// wait for build and run the binary with given params
	go r.Init(params)

	// build given package
	go build(w, r, params)

	// force update for initial package build
	go w.ForceUpdate()

	// listen for further changes
	go w.ListenChanges()

	<-done
}
