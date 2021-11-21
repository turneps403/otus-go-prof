package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	proc := exec.Command(cmd[0], cmd[1:]...)

	proc.Env = make([]string, len(os.Environ()))
	copy(proc.Env, os.Environ())
	for k, v := range env {
		if v.NeedRemove {
			proc.Env = append(proc.Env, fmt.Sprintf("%v=", k))
		} else {
			proc.Env = append(proc.Env, fmt.Sprintf("%v=%v", k, v))
		}
	}

	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	err := proc.Run()
	if err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}
