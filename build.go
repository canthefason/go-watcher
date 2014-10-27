package main

import "log"

func build(w *Watcher, r *Runner) {
	for {
		w.Wait()

		log.Println("Building...")

		// TODO build path must be optional
		cmd, err := runCommand("go", "build", "-o", binaryName)
		if err != nil {
			log.Fatalf("Could not run 'go build' command: %s", err)
			continue
		}

		if err := cmd.Wait(); err != nil {
			log.Println("Waiting for file change...")
			continue
		}

		// when binary is successfully updated, kill the old running process
		r.Kill()

		// and start the new process
		r.Run()
	}
}
