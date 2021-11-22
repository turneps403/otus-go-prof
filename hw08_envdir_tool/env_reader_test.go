package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func tmpDir() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), fmt.Sprintf("%v%v", hex.EncodeToString(randBytes), os.Getegid()))
}

func createFileWithContent(dir, name, content string) (string, error) {
	srcFileName := filepath.Join(dir, name)
	srcFile, err := os.OpenFile(srcFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	if len(content) > 0 {
		fmt.Fprint(srcFile, content)
	}
	srcFile.Close()
	return srcFileName, nil
}

func TestReadDir(t *testing.T) {
	t.Log("run TestReadDir")

	dir := tmpDir()
	err := os.Mkdir(dir, os.ModePerm)
	require.NoError(t, err)
	defer os.Remove(dir)

	fmt.Printf("dir: %v\n", dir)

	t.Run("Filename starts with a number", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "9FOO", "123")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, len(envMap), 0)
	})

	t.Run("Normal filename", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO", "123")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, val.Value, "123")
		require.Equal(t, val.NeedRemove, false)
	})

	t.Run("Normal filename with underscore", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO_BAR", "123")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO_BAR"]
		require.Equal(t, ok, true)
		require.Equal(t, val.Value, "123")
		require.Equal(t, val.NeedRemove, false)
	})

	t.Run("File no content", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO", "")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, val.Value, "")
		require.Equal(t, val.NeedRemove, true)
	})

	t.Run("Multiline content", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO", "foo\nbar\nbaz")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, val.Value, "foo")
		require.Equal(t, val.NeedRemove, false)
	})

	t.Run("Content with tab", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO", "\tfoo")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, "\"\\tfoo\"", val.Value)
		require.Equal(t, val.NeedRemove, false)
	})

	t.Run("Content with indent", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "FOO", "    foo")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, "    foo", val.Value)
		require.Equal(t, val.NeedRemove, false)
	})

	t.Run("Filename contains lowercase letters", func(t *testing.T) {
		filePath, err := createFileWithContent(dir, "Foo", "123")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, len(envMap), 0)
	})

	t.Run("Content with nullbyte", func(t *testing.T) {
		// normal file
		filePath, err := createFileWithContent(dir, "FOO", "fo"+string('\x00')+"o")
		require.NoError(t, err)
		defer os.Remove(filePath)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		val, ok := envMap["FOO"]
		require.Equal(t, ok, true)
		require.Equal(t, "\"fo\\x00o\"", val.Value)
		require.Equal(t, val.NeedRemove, false)
	})
}
