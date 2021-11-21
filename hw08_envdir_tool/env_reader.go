package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
	"sync"
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

func isValidEnvName(name string) bool {
	for _, c := range name {
		if ('A' > c || c > 'Z') && c != '_' {
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
			if file.Size() == 0 {
				myEnv.add(file.Name(), "", true)
				return
			}
			content, err := ioutil.ReadFile(file.Name())
			if err != nil {
				log.Fatal(err)
				return
			}
			myEnv.add(file.Name(), strings.TrimSpace(string(content)), false)
		}(file)
	}
	wg.Wait()

	return myEnv.env, nil
}
