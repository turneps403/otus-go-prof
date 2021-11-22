package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	proc := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	proc.Env = make([]string, len(os.Environ()))
	copy(proc.Env, os.Environ())
	for k, v := range env {
		if v.NeedRemove {
			proc.Env = append(proc.Env, fmt.Sprintf("%v=", k))
		} else {
			proc.Env = append(proc.Env, fmt.Sprintf("%v=%v", k, v.Value))
		}
	}

	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	if err := proc.Start(); err != nil {
		log.Fatalf("proc.Start: %v", err)
	}

	if err := proc.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok { //nolint:errorlint
			// see details https://stackoverflow.com/questions/10385551/get-exit-code-go
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			log.Fatalf("proc.Wait: %v", err)
		}
	}

	return 0
}
