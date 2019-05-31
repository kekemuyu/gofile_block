package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
)

type Server struct {
	Conn net.Conn
	Flag bool
}

var DefaultServer = Server{}

func (s *Server) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("有一个客户端上线：", conn.RemoteAddr().String())

		DefaultServer.Conn = conn
		DefaultServer.Flag = true
	}

}

func (s *Server) Read() bytes.Buffer {
	if s.Flag {
		buf, err := ioutil.ReadAll(s.Conn)
		if err != nil {
			panic(err)
		}

		var bb bytes.Buffer
		_, err = bb.Write(buf)
		if err != nil {
			panic(err)
		}
		return bb
	} else {
		return bytes.Buffer{}
	}

}

func (s *Server) Write(bytes.Buffer) {}
