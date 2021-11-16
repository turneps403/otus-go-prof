package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	PagePerCopy = 1
)

var (
	ErrBadArguments = errors.New("Bad arguments")
	ErrBadFile      = errors.New("Bad file")
	ErrBadCopy      = errors.New("Copy problem")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return errors.Wrap(ErrBadArguments, "'limit' shouldn't be less than zero")
	}
	// open a destination file
	_, err := os.Stat(toPath)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(ErrBadFile, fmt.Sprintf("destination file %s already exists", toPath))
	}
	to, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(ErrBadFile, err.Error())
	}
	defer to.Close()

	// checking permissions of source file
	if info, err := os.Stat(fromPath); err == nil {
		if !info.Mode().IsRegular() {
			return errors.Wrap(ErrBadFile, fmt.Sprintf("%s is not a regular file", fromPath))
		}
		if offset < 0 {
			offset = info.Size() + offset
		}
		if offset < 0 || offset > info.Size() {
			return errors.Wrap(ErrBadArguments, fmt.Sprintf("'offset' isn't suitable for a file with size %d\n", info.Size()))
		}
		if offset == info.Size() {
			return nil
		}
		if limit == 0 {
			limit = info.Size() - offset
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(ErrBadFile, fmt.Sprintf("path '%s' isn't a file: %v", fromPath, err))
	} else {
		// panic(err)
		return errors.Wrap(ErrBadFile, err.Error())
	}

	// open a source file
	from, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
	if err != nil {
		return errors.Wrap(ErrBadFile, err.Error())
	}
	defer to.Close()

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.Wrap(ErrBadFile, fmt.Sprintf("seek problem: %s", err.Error()))
	}

	copyLimit := int64(PagePerCopy * os.Getpagesize())
	for limit > 0 {
		if copyLimit > limit {
			copyLimit = limit
		}
		if wCnt, err := io.CopyN(to, from, copyLimit); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.Wrap(ErrBadCopy, err.Error())
		} else {
			limit -= wCnt
		}
	}

	return nil
}
