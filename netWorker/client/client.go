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

func (c *Client) Write(bs []byte) {
	_, err := c.Conn.Write(bs)
	if err != nil {
		panic(err)
	}
}

func (c *Client) Read() bytes.Buffer {
	return bytes.Buffer{}
}
