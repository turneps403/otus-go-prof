package hw02unpackstring

import (
	"fmt"
	"strings"

	"github.com/turneps403/otus-go-prof/hw02_unpack_string/myreader"
)

type UnpackError struct {
	reason string
	err    error
}

func (e *UnpackError) Error() string {
	return fmt.Sprintf("reason: %s, because of %v", e.reason, e.err)
}

func Unpack(s string) (string, error) {
	var sb strings.Builder
	r := myreader.NewMyReader(s)
	for r.HasNext() {
		ru, rep, err := r.Next()
		if err != nil {
			if err, ok := err.(*myreader.ReaderError); ok {
				return "", &UnpackError{reason: "reader cant get next character", err: err}
			} else {
				return "", err
			}
		}
		fmt.Fprint(&sb, strings.Repeat(string(ru), rep))
	}
	return sb.String(), nil
}
