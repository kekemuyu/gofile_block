package main

import (
	"flag"
	"fmt"
	"io"

	//	"time"

	"gofile_block/internal"
	"gofile_block/networker/client"
	"gofile_block/networker/server"
	//	"gofile_block/pipe"
	//	"gofile_block/serialworker"
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

	if *hostname != "" {
		fmt.Println("send file")
		c := client.New(*hostname)

		if *file == "" {
			fmt.Println("input the file name with -f ")
			return
		}

		internal.Defaultbuffer.OpenFile(*file)
		internal.Defaultbuffer.GetFileInfo()
		internal.Defaultbuffer.Fileinfo.Off = 0
		internal.Defaultbuffer.Fileinfo.Start = 0xAA

		for {

			bs, err := internal.Defaultbuffer.ReadBlock()
			if (err != nil) && (err != io.EOF) {
				panic(err)
			}

			if len(bs) > 0 {
				buf := make([]byte, 3)
				buf[0] = byte(len(bs) & 0xFF)
				buf[1] = byte((len(bs) >> 8) & 0xFF)
				buf[2] = 0
				if internal.Defaultbuffer.Fileinfo.Start == 0xAA {
					buf[2] = 0xAA
					internal.Defaultbuffer.Fileinfo.Start = 0
				}
				buf = append(buf, bs...)
				//				fmt.Println(buf)
				c.Write(buf)

			}
			//			fmt.Println(internal.Defaultbuffer.Fileinfo.Header.Size, internal.Defaultbuffer.Fileinfo.Off)
			if internal.Defaultbuffer.Fileinfo.Header.Size == internal.Defaultbuffer.Fileinfo.Off {
				fmt.Println("send complite!")
				internal.Defaultbuffer.Fileinfo.Handle.Close()
				c.Conn.Close()
				return
			}

		}

	} else if *port != "" {

		go server.DefaultServer.Run(*port)
		for {
			if server.DefaultServer.Flag != true {
				continue
			}

			buf := make([]byte, 3)

			n, err := server.DefaultServer.Conn.Read(buf)
			bsLens := int(buf[1])<<8 + int(buf[0])
			if (n == 0) || (err != nil) {
				fmt.Println("文件接收完毕")
				internal.Defaultbuffer.Fileinfo.Handle.Close()
				server.DefaultServer.Conn.Close()
				server.DefaultServer.Flag = false
				continue
			}
			if buf[2] == 0xAA {
				internal.Defaultbuffer.Fileinfo.Start = 0xAA
			} else {
				internal.Defaultbuffer.Fileinfo.Start = 0
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
		//		fmt.Println("opened com port is:", *com)
		//		w := serialworker.New(*com)
		//		ctrl = w
		//		if *mode == "s" {
		//			if *file == "" {
		//				fmt.Println("input the file name with -f ")
		//				return
		//			}
		//			bbytes := internal.Defaultbuffer.GetBytesbuffer(*file)
		//			ctrl.Write(bbytes)
		//		} else if *mode == "r" {
		//			for {
		//				bb := ctrl.Read()
		//				internal.Defaultbuffer.PutBytesbufferToFile(bb.Bytes())

		//			}
		//		} else {
		//			fmt.Println("input the serial mode  with -r or -s ")
		//		}

	} else {
		fmt.Println("Please input gofile -h for help")
	}
}
