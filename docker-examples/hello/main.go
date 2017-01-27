package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":7000",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "watcher is running")
			},
		),
	)
}
