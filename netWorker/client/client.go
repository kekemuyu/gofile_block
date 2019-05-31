package client

import (
	"bytes"
	"net"
)

type Client struct {
	Conn net.Conn
}

func New(hostname string) *Client {
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		panic(conn)
	}

	return &Client{
		Conn: conn,
	}
}

func (c *Client) Write(bb bytes.Buffer) {
	_, err := c.Conn.Write(bb.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *Client) Read() bytes.Buffer {
	return bytes.Buffer{}
}
