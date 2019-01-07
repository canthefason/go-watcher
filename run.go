// Package watcher is a command line tool inspired by fresh (https://github.com/pilu/fresh) and used
// for watching .go file changes, and restarting the app in case of an update/delete/add operation.
// After you installed it, you can run your apps with their default parameters as:
// watcher -c config -p 7000 -h localhost
package watcher

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

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
			color.Red(fmt.Sprintf("Could not run the go binary: %s \n", err))
			r.kill(cmd)

			continue
		}

		r.cmd = cmd
		removeFile(fileName)

		go func(cmd *exec.Cmd) {
			if err := cmd.Wait(); err != nil {
				color.Red(fmt.Sprintf("process interrupted: %s \n", err))
				r.kill(cmd)
			}
		}(r.cmd)
	}
}

// Restart kills the process, removes the old binary and
// restarts the new process
func (r *Runner) restart(fileName string) {
	r.kill(r.cmd)

	r.start <- fileName
}

// Signal the process to shutdown, if it doesn't, kill it!
func (r *Runner) kill(cmd *exec.Cmd) {
	if cmd != nil {
		cmd.Process.Signal(syscall.SIGINT)

		didExit := make(chan struct{})
		go func() {
			select {
			case <-didExit:
			case <-time.After(5 * time.Second):
				cmd.Process.Kill()
			}
		}()

		state, err := cmd.Process.Wait()
		if err == nil && state.Exited() {
			close(didExit)
		}
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
