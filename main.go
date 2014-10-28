package main

const binaryName = "goldorf-main"

func main() {
	params := prepareArgs()

	w := MustRegisterWatcher(params)
	defer w.Close()

	done := make(chan struct{})

	r := NewRunner()

	// wait for build and run the binary with given params
	go r.Init(params)

	// first build given package
	go build(w, r, params)

	// force update for initial package build
	go w.ForceUpdate()

	// listen for further changes
	go w.ListenChanges()

	<-done
}
