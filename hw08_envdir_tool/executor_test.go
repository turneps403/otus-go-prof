package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmdEnv(t *testing.T) {
	t.Log("run env")
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	exCode := RunCmd([]string{"env"}, make(Environment))
	require.Equal(t, 0, exCode)
}

func TestRunCmdLsLa(t *testing.T) {
	t.Log("run ls -la")
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	exCode := RunCmd([]string{"ls", "-la"}, make(Environment))
	require.Equal(t, 0, exCode)
}

func TestRunCmdFailArg(t *testing.T) {
	t.Log("run ls -KBD")
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)
	exCode := RunCmd([]string{"ls", "-KBD"}, make(Environment))
	require.Equal(t, 1, exCode)
}

func TestRunCmdPanic(t *testing.T) {
	t.Log("run broken env")

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
	exCode := RunCmd([]string{"env"}, e)
	require.Equal(t, -1, exCode) // never happens
}
