package main

import (
	"log"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	update  chan struct{}
}

func MustRegisterWatcher() *Watcher {
	w := &Watcher{
		update: make(chan struct{}),
	}

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not register watcher: %s", err)
	}

	// TODO configurable input path
	// add watched paths
	err = w.watcher.Add("./")
	if err != nil {
		log.Fatalf("Could not register watcher: %s", err)
	}

	// send update signal for initial package build
	go func() {
		w.update <- struct{}{}
	}()

	return w
}

// ListenChanges listens file updates, and sends signal to
// update channel when go files are updated
func (w *Watcher) ListenChanges() {
	for {
		select {
		case event := <-w.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {

				ext := filepath.Ext(event.Name)
				if ext == ".go" || ext == ".tmpl" {
					w.update <- struct{}{}
				}

			}
		case err := <-w.watcher.Errors:
			log.Fatalf("Watcher error: %s", err)
		}
	}
}

func (w *Watcher) Close() {
	w.watcher.Close()
}

func (w *Watcher) Wait() {
	<-w.update
}
