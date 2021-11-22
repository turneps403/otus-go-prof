package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func showDoc(err error) {
	log.Fatalf("%v\n-> use: %v /some/dir prog arg1 ...", err.Error(), filepath.Base(os.Args[0]))
}

func main() {
	if len(os.Args) < 3 {
		showDoc(fmt.Errorf("not enough arguments"))
	}
	envMap, err := ReadDir(os.Args[1])
	if err != nil {
		showDoc(err)
	}
	exCode := RunCmd(os.Args[2:], envMap)
	os.Exit(exCode)
}
