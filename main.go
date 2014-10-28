package main

const binaryName = "goldorf-main"

func main() {
	params := prepareArgs()

	w := MustRegisterWatcher(rootpackage)
	defer w.Close()

	done := make(chan struct{})

	r := NewRunner()

	// first build given package
	go build(w, r, params)

	// run the binary with given arguments
	go r.Init(args...)
	// force update for initial package build
	go w.ForceUpdate()

	// listen for further changes
	go w.ListenChanges()

	<-done
}
