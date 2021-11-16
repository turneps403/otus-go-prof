package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyArgsError(t *testing.T) {
	t.Log("run TestCopyArgsError")

	// source and destination empty
	err := Copy("", "", 0, 0)
	require.Error(t, err)

	// emty source
	err = Copy("", "/foo/bar", 0, 0)
	require.Error(t, err)

	tmpFile, err := ioutil.TempFile("/tmp", "test.*.xxx")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// special file
	err = Copy("/dev/urandom", tmpFile.Name(), 0, 0)
	require.Error(t, err)

	// tmpFile.Name() exists
	err = Copy(tmpFile.Name(), tmpFile.Name(), 0, 0)
	require.Error(t, err)

	// limit less than 0
	err = Copy(tmpFile.Name(), tmpFile.Name(), 0, -5)
	require.Error(t, err)
}

func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return filepath.Join(os.TempDir(), prefix+"fail_foo_bar"+suffix)
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
	// return filepath.Join("/tmp", prefix+hex.EncodeToString(randBytes)+suffix)
}

func TestCopySuccess(t *testing.T) {
	t.Log("run TestCopySuccess")

	srcFileName := tempFileName("foo", ".src")
	srcFile, err := os.OpenFile(srcFileName, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	defer func() {
		srcFile.Close()
		os.Remove(srcFileName)
	}()

	testContent := "1234567890"
	fmt.Fprint(srcFile, testContent)
	srcFile.Close()

	dstFileName := tempFileName("foo", ".dst")
	defer os.Remove(dstFileName)

	// full copy
	err = Copy(srcFileName, dstFileName, 0, 0)
	require.NoError(t, err)
	dstContent, _ := ioutil.ReadFile(dstFileName)
	require.Equal(t, testContent, string(dstContent))

	// head copy
	os.Remove(dstFileName)
	err = Copy(srcFileName, dstFileName, 0, 4)
	require.NoError(t, err)
	dstContent, _ = ioutil.ReadFile(dstFileName)
	require.Equal(t, testContent[:4], string(dstContent))

	// tail copy
	os.Remove(dstFileName)
	err = Copy(srcFileName, dstFileName, -4, 0)
	require.NoError(t, err)
	dstContent, _ = ioutil.ReadFile(dstFileName)
	require.Equal(t, testContent[len(testContent)-4:], string(dstContent))

	// part copy
	os.Remove(dstFileName)
	err = Copy(srcFileName, dstFileName, -4, 2)
	require.NoError(t, err)
	dstContent, _ = ioutil.ReadFile(dstFileName)
	require.Equal(t, testContent[len(testContent)-4:len(testContent)-4+2], string(dstContent))
}
