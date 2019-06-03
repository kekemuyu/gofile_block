package server

import (
	"bytes"
	"fmt"
	"net"
)

type Server struct {
	Conn net.Conn
	Flag bool
}

var DefaultServer Server

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

		s.Conn = conn
		s.Flag = true
	}

}

func (s *Server) Write(bytes.Buffer) {}
