package myreader

import (
	"fmt"
)

// https://www.youtube.com/watch?v=oIxXp0OgK_0

type ReaderError struct {
	reason string
	err    error
}

func (e *ReaderError) Error() string {
	return fmt.Sprintf("reason: %s", e.reason)
}

func (e *ReaderError) Unwrap() error {
	return e.err
}
