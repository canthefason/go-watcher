// Package watcher is a command line tool inspired by fresh (https://github.com/pilu/fresh) and used
// for watching .go file changes, and restarting the app in case of an update/delete/add operation.
// After you installed it, you can run your apps with their default parameters as:
// watcher -c config -p 7000 -h localhost
package watcher

import (
	"log"
	"os/exec"

	"github.com/fatih/color"
)

// Runner listens for the change events and depending on that kills
// the obsolete process, and runs a new one
type Runner struct {
	start chan string
	done  chan struct{}
	cmd   *exec.Cmd
}

// NewRunner creates a new Runner instance and returns its pointer
func NewRunner() *Runner {
	return &Runner{
		start: make(chan string),
		done:  make(chan struct{}),
	}
}

// Run initializes runner with given parameters.
func (r *Runner) Run(p *Params) {
	for fileName := range r.start {

		color.Green("Running %s...\n", p.Get("run"))

		cmd, err := runCommand(fileName, p.Package...)
		if err != nil {
			log.Printf("Could not run the go binary: %s \n", err)
			r.kill(cmd)

			continue
		}

		r.cmd = cmd

		go func(cmd *exec.Cmd, fileName string) {
			if err := cmd.Wait(); err != nil {
				log.Printf("process interrupted: %s \n", err)
				r.kill(cmd)
			}

			removeFile(fileName)
		}(r.cmd, fileName)
	}
}

// Restart kills the process, removes the old binary and
// restarts the new process
func (r *Runner) restart(fileName string) {
	r.kill(r.cmd)

	r.start <- fileName
}

func (r *Runner) kill(cmd *exec.Cmd) {
	if cmd != nil {
		cmd.Process.Kill()
	}
}

func (r *Runner) Close() {
	close(r.start)
	r.kill(r.cmd)
	close(r.done)
}

func (r *Runner) Wait() {
	<-r.done
}
