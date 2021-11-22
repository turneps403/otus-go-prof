package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunGit(t *testing.T) {
	t.Log("run git --version")
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	exCode := RunCmd([]string{"git", "--version"}, make(Environment))
	require.Equal(t, 0, exCode)
}

func TestRunCmdFailArg(t *testing.T) {
	t.Log("run git -KBD")
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)
	exCode := RunCmd([]string{"git", "--KBD"}, make(Environment))
	require.Equal(t, true, exCode == 129)
}

func TestRunCmdPanic(t *testing.T) {
	t.Log("run panic")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	e := make(Environment)
	e["FOO"] = EnvValue{Value: "foo" + string('\x00') + "bar"}
	exCode := RunCmd([]string{"git", "--version"}, e)
	require.Equal(t, -1, exCode) // never happens
}
