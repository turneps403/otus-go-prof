package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/turneps403/otus-go-prof/hw11_telnet_client/logger"
)

// Parsing incomin args: prog --timeout=10s host port.
func parseArgs() (host string, port int, timeout time.Duration) {
	zap := logger.Zap()
	refTimeout := flag.Duration("timeout", 10*time.Second, "definition of time for awaiting new connection")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Please use: %s telnet host port [--timeout=2s]\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	host = flag.Args()[0]
	if p, err := strconv.Atoi(flag.Args()[1]); err != nil {
		zap.Fatal(err)
	} else {
		if p < 1 || p > 65_535 {
			zap.Fatal("Port should be in range among 1 and 65 535")
		}
		port = p
	}
	timeout = *refTimeout
	return
}

func main() {
	defer logger.Finalize()
	zap := logger.Zap()

	host, port, timeout := parseArgs()

	cl := NewTelnetClient(
		net.JoinHostPort(host, strconv.Itoa(port)),
		timeout,
		os.Stdin,
		os.Stdout)

	if err := cl.Connect(); err != nil {
		zap.Fatal(err)
	}
	defer cl.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Buffer was add because of:
	// "misuse of unbuffered os.Signal channel as argument to signal.Notify".
	signalsChan := make(chan os.Signal, 10)
	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)

	// Sending smth
	go func() {
		zap := logger.Zap()
		cl.Send()
		err := cl.Send()
		if err != nil {
			zap.Error(err)
		}
		cancel()
	}()

	// Getting smth
	go func() {
		zap := logger.Zap()
		err := cl.Receive()
		if err != nil {
			zap.Error(err)
		}
		cancel()
	}()

	// Awaiting signals and graceful shutdown
	go func() {
		select {
		case <-ctx.Done():
		case <-signalsChan:
			zap := logger.Zap()
			zap.Info("Connection was closed by peer")
			cancel()
		}
	}()

	<-ctx.Done()
}
