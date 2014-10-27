package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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
	w.watchFolders()

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

// watchFolders recursively adds folders that will be watched for changes,
// starting from the working directory
func (w *Watcher) watchFolders() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get current working directory: %s", err)
	}

	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		// skip files
		if !info.IsDir() {
			return nil
		}

		// skip hidden folders
		if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		w.addFolder(path)

		return err
	})
}

// addFolder adds given folder name to the watched folders, and starts
// watching it for further changes
func (w *Watcher) addFolder(name string) {
	if err := w.watcher.Add(name); err != nil {
		log.Fatalf("Could not watch folder: %s", err)
	}
}
