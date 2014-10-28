package main

import (
	"io"
	"log"
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

func prepareArgs() *Params {

	params := NewParams()

	args := os.Args

	// remove command
	args = args[1:len(args)]

	for i := 0; i < len(args); i++ {
		arg := args[i]
		arg = stripDash(arg)

		if arg == "run" || arg == "watch" {
			if len(args) <= i+1 {
				log.Fatalf("missing parameter value: %s", arg)
			}

			params.System[arg] = args[i+1]
			i++
			continue
		}

		params.Package = append(params.Package, args[i])
	}

	params.CloneRun()

	return params
}

func stripDash(arg string) string {
	if len(arg) > 2 {
		if arg[1] == '-' {
			return arg[2:]
		} else if arg[0] == '-' {
			return arg[1:]
		}
	}

	return arg
}
