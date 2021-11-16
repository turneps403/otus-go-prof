package main

import (
	"os"
	"os/signal"
	"syscall"
)

type ctrlC struct {
	cb   func()
	done chan bool
	c    chan os.Signal
}

func NewCtrlC(cb func()) *ctrlC {
	ret := &ctrlC{
		c:    make(chan os.Signal, 10),
		done: make(chan bool),
		cb:   cb,
	}
	signal.Notify(ret.c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-ret.c:
			cb()
			os.Exit(1)
		case <-ret.done:
		}
	}()

	return ret
}

func (ctrlc *ctrlC) Done() {
	ctrlc.done <- true
	close(ctrlc.done)
	close(ctrlc.c)
}
