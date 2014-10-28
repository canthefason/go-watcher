package main

import (
	"io"
	"os"
	"os/exec"
)

type Params struct {
	Package []string
	System  map[string]string
}

func NewParams() *Params {
	return &Params{
		Package: make([]string, 0),
		System:  make(map[string]string),
	}
}

func (p *Params) Get(name string) string {
	return p.System[name]
}

func (p *Params) CloneRun() {
	if p.System["watch"] == "" && p.System["run"] != "" {
		p.System["watch"] = p.System["run"]
	}
}

func runCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, err
	}

	if err := cmd.Start(); err != nil {
		return cmd, err
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd, nil
}
