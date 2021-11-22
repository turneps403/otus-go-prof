package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type syncEnvironment struct {
	mu  sync.Mutex
	env Environment
}

func (c *syncEnvironment) add(name, value string, remove bool) {
	c.mu.Lock()
	c.env[name] = EnvValue{Value: value, NeedRemove: remove}
	c.mu.Unlock()
}

func chompLine(s string) string {
	// to simplify things, we read whole content of file
	// and chomp one line from the top.
	// not a production version, just a concept
	hasNonGraphic, isOnlySpace := false, true
	sb := strings.Builder{}
	for _, c := range s {
		if c == '\n' {
			break
		}
		if !hasNonGraphic && !unicode.IsGraphic(c) {
			hasNonGraphic = true
		}
		if isOnlySpace && !unicode.IsSpace(c) {
			isOnlySpace = false
		}
		sb.WriteRune(c)
	}
	if hasNonGraphic || isOnlySpace {
		return fmt.Sprintf("%q", sb.String())
	}
	return sb.String()
}

func isValidEnvName(name string) bool {
	if name == "NUL" {
		return false
	}
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}
	for _, c := range name {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	myEnv := &syncEnvironment{env: make(Environment)}
	var wg sync.WaitGroup
	for _, file := range files {
		if !isValidEnvName(file.Name()) {
			continue
		}
		wg.Add(1)
		go func(finfo fs.FileInfo) {
			defer wg.Done()
			if finfo.Size() == 0 {
				myEnv.add(finfo.Name(), "", true)
				return
			}
			content, err := ioutil.ReadFile(filepath.Join(dir, finfo.Name()))
			if err != nil {
				log.Fatal(err)
				return
			}
			myEnv.add(finfo.Name(), chompLine(string(content)), false)
		}(file)
	}
	wg.Wait()

	return myEnv.env, nil
}
