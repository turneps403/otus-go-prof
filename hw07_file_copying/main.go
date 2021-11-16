package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	from, to = normalizePath(from), normalizePath(to)
	log.Printf("source file path '%s' and destination '%s'", from, to)

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatal(err)
	}
}

func normalizePath(path string) string {
	if len(path) > 1 && path[:2] == "./" {
		workDir, _ := os.Getwd()
		path = fmt.Sprintf("%s/%s", workDir, path[2:])
	}
	if len(path) > 1 && path[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		path = fmt.Sprintf("%s/%s", homeDir, path[2:])
	}
	if len(path) > 0 && path[:1] != "/" {
		workDir, _ := os.Getwd()
		path = fmt.Sprintf("%s/%s", workDir, path)
	}
	path, _ = filepath.Abs(path)
	return path
}

func usage() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		file = filepath.Base(file)
	} else {
		file = " THIS_BIN"
	}
	fmt.Fprintf(os.Stderr, "usage: go run %s -from=xxx -to=xxx -limit=xxx -offset=xxx\n", file)
	flag.PrintDefaults()
	os.Exit(2)
}
