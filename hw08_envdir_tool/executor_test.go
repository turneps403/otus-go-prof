package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunGit(t *testing.T) {
	t.Log("run git --version")
	exCode := RunCmd([]string{"git", "--version"}, make(Environment))
	require.Equal(t, 0, exCode)
}

func TestRunCmdFailArg(t *testing.T) {
	t.Log("run git -KBD")
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

	e := make(Environment)
	e["FOO"] = EnvValue{Value: "foo" + string('\x00') + "bar"}
	exCode := RunCmd([]string{"git", "--version"}, e)
	require.Equal(t, -1, exCode) // never happens
}
