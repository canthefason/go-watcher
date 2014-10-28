package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	rootdir string
	watcher *fsnotify.Watcher
	update  chan struct{}
}

var ErrPathNotSet = errors.New("gopath not set")

func MustRegisterWatcher(params *Params) *Watcher {

	w := &Watcher{
		update:  make(chan struct{}),
		rootdir: params.Get("watch"),
	}

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not register watcher: %s", err)
	}

	// TODO configurable input path
	// add watched paths
	w.watchFolders()

	return w
}

// ListenChanges listens file updates, and sends signal to
// update channel when go files are updated
func (w *Watcher) ListenChanges() {
	for {
		select {
		case event := <-w.watcher.Events:
			if event.Op&fsnotify.Chmod != fsnotify.Chmod {
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

// send update signal for initial package build
func (w *Watcher) ForceUpdate() {
	w.update <- struct{}{}
}

// watchFolders recursively adds folders that will be watched for changes,
// starting from the working directory
func (w *Watcher) watchFolders() {
	wd, err := w.prepareRootDir()

	if err != nil {
		log.Fatalf("Could not get root working directory: %s", err)
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

func (w *Watcher) prepareRootDir() (string, error) {
	if w.rootdir == "" {
		return os.Getwd()
	}

	path := os.Getenv("GOPATH")
	if path == "" {
		return "", ErrPathNotSet
	}

	root := fmt.Sprintf("%s/src/%s", path, w.rootdir)

	return root, nil
}
