package main

import (
	"flag"
	"fmt"
	"io"

	//	"time"

	"gofile_block/internal"
	"gofile_block/networker/client"
	"gofile_block/networker/server"
	"gofile_block/pipe"
	"gofile_block/serialworker"
)

var file = flag.String("f", "test.txt", "input the file name")
var hostname = flag.String("s", "", "input server ip")
var port = flag.String("p", "", "input server port")
var com = flag.String("c", "", "input com port")
var mode = flag.String("m", "", "input r or s for recieve or send ")

func init() {
	flag.Parse()
}

func main() {
	var ctrl pipe.Control
	if *hostname != "" {
		fmt.Println("send file")
		c := client.New(*hostname)
		ctrl = c
		if *file == "" {
			fmt.Println("input the file name with -f ")
			return
		}

		internal.Defaultbuffer.Create(*file)

		for {
			bbytes, err := internal.Defaultbuffer.ReadBlock()
			if (err != nil) && (err != io.EOF) {
				panic(err)
			}
			//			fmt.Println(bbytes)
			if len(bbytes.Bytes()) > 0 {
				//				fmt.Println(bbytes)
				ctrl.Write(bbytes)
				//				time.Sleep(time.Microsecond * 100)
			}
			if err == io.EOF {
				fmt.Println("send file complite!")
				for {
				}
			}
		}

	} else if *port != "" {

		go server.DefaultServer.Run(*port)
		for {
			if server.DefaultServer.Flag != true {
				continue
			}
			//			fmt.Println(server.DefaultServer.Flag)
			buf := make([]byte, 2)

			n, err := server.DefaultServer.Conn.Read(buf)
			bsLens := int(buf[1])<<8 + int(buf[0])
			if bsLens == 0 {
				fmt.Println("读取到长度为0")
				for {
				}
			}
			buf = make([]byte, bsLens)

			n, err = server.DefaultServer.Conn.Read(buf)
			if err != nil {
				panic(err)
			}
			//			fmt.Println(buf[:n])
			internal.Defaultbuffer.WriteBlock(buf[:n])

		}
	} else if (*com != "") && (*mode != "") {
		fmt.Println("opened com port is:", *com)
		w := serialworker.New(*com)
		ctrl = w
		if *mode == "s" {
			if *file == "" {
				fmt.Println("input the file name with -f ")
				return
			}
			bbytes := internal.Defaultbuffer.GetBytesbuffer(*file)
			ctrl.Write(bbytes)
		} else if *mode == "r" {
			for {
				bb := ctrl.Read()
				internal.Defaultbuffer.PutBytesbufferToFile(bb.Bytes())

			}
		} else {
			fmt.Println("input the serial mode  with -r or -s ")
		}

	} else {
		fmt.Println("Please input gofile -h for help")
	}
}
