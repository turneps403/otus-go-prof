package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TCPClient struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *TCPClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}

func (c *TCPClient) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return err
	}
	return nil
}

func (c *TCPClient) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPClient{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// P.S. Author's solution takes no more than 50 lines.
