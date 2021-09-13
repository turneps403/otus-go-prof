package myreader

import "fmt"

type ReaderError struct {
	reason string
}

func (e *ReaderError) Error() string {
	return fmt.Sprintf("reason: %s", e.reason)
}
