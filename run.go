package main

import (
	"fmt"
	"log"
)

type Runner struct {
	running bool
	start   chan struct{}
	kill    chan struct{}
}

func NewRunner() *Runner {
	return &Runner{
		running: false,
		start:   make(chan struct{}),
		kill:    make(chan struct{}),
	}
}

func (r *Runner) Init(args ...string) {

	for {
		<-r.start

		log.Println("Running...")

		cmd, err := runCommand(prepareBinaryName(binaryName), args...)
		if err != nil {
			log.Println("Could not run the go binary: %s", err)
			continue
		}

		go func() {
			r.running = true
			cmd.Wait()
		}()

		go func() {
			<-r.kill
			pid := cmd.Process.Pid
			log.Printf("Killing PID %d \n", pid)
			cmd.Process.Kill()
		}()

	}
}

func prepareBinaryName(name string) string {
	return fmt.Sprintf("./%s", name)
}

func (r *Runner) Run() {
	r.start <- struct{}{}
}

func (r *Runner) Kill() {
	if r.running {
		r.kill <- struct{}{}
	}
}
