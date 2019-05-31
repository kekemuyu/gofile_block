package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"

	"os"
)

const Blocksize = 1024

type File struct {
	Handle *os.File //文件
	Start  byte     //开始标志
	Name   string   //文件名称
	Size   int64    //文件大小
	Off    int64    //文件位置偏移
}

type Buffer struct {
	Fileinfo File
}

var Defaultbuffer Buffer

func (b *Buffer) OpenFile(name string) {
	fin, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	b.Fileinfo.Handle = fin
}

func (b *Buffer) GetFileInfo(name string) {

	fileInfo, err := fin.Stat()
	if err != nil {
		panic(err)
	}

	fileObj := File{
		Handle: b.Fileinfo.Handle,
		Start:  0xAA,
		Name:   fileInfo.Name(),
		Size:   fileInfo.Size(),
		Off:    0,
	}

	b.Fileinfo = fileObj
}

func (b *Buffer) ReadBlock() (bytes.Buffer, int, error) {
	var outbb bytes.Buffer
	var bs []byte
	var err error

	bsize := Blocksize
	if b.Fileinfo.Size <= 0 {
		return bytes.Buffer{}, 0, errors.New("文件大小为0")
	} else if b.Fileinfo.Size < Blocksize {
		bsize = int(b.Fileinfo.Size)
	} else if (b.Fileinfo.Size - b.Fileinfo.Size.Off) < Blocksize {
		bsize = int(b.Fileinfo.Size - Defaultbuffer.FileInfo.Off)
	}

	if bsize == 0 {
		bsize = 1
	}
	bs = make([]byte, bsize)
	n, err := b.Fileinfo.Handle.ReadAt(bs, Defaultbuffer.FileInfo.Off)
	if n <= 0 {
		fmt.Println(err)
		return bytes.Buffer{}, 0, err
	}

	if b.Fileinfo.Start == 0xAA { //如果是开始发送第一包
		b.Fileinfo.Start = 0x00

		bs := make()
	}
	outbs, err := json.Marshal(Defaultbuffer.FileInfo)
	if err != nil {
		fmt.Println("ReadBlock marshal error", err)
		return bytes.Buffer{}, err
	}
	bsLen := len(outbs)
	outbb.WriteByte(byte(bsLen & 0xFF))
	outbb.WriteByte(byte((bsLen >> 8) & 0xFF))
	if _, err = outbb.Write(outbs); err != nil {
		fmt.Println("ReadBlock writebuffer error", err)
		return bytes.Buffer{}, err
	}
	fmt.Println(Defaultbuffer.FileInfo.Off)
	Defaultbuffer.FileInfo.Off = Defaultbuffer.FileInfo.Off + int64(n)

	return outbb, err

}

func (b *Buffer) WriteBlock(bs []byte) {
	var err error

	fileinfo := File{
		Data: make([]byte, Blocksize),
	}
	err = json.Unmarshal(bs, &fileinfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	Defaultbuffer.FileInfo = fileinfo

	if Defaultbuffer.Filein == nil {
		if Defaultbuffer.Filein, err = os.Create(fileinfo.Name); err != nil {
			fmt.Println("create file error", err)
			return
		}
		if Defaultbuffer.Filein, err = os.Create(fileinfo.Name); err != nil {
			panic(err)
		}
		Defaultbuffer.FileInfo.Off = 0
	}
	Defaultbuffer.FileInfo.Data = fileinfo.Data
	bsize := Blocksize
	if len(Defaultbuffer.FileInfo.Data) < bsize {
		bsize = len(Defaultbuffer.FileInfo.Data)
	}

	var wlen int
	//	fmt.Println(Defaultbuffer.FileInfo)
	if wlen, err = Defaultbuffer.Filein.WriteAt(Defaultbuffer.FileInfo.Data, Defaultbuffer.FileInfo.Off); err != nil {
		panic(err)
		return
	}
	fmt.Println(Defaultbuffer.FileInfo.Off)
	Defaultbuffer.FileInfo.Off = Defaultbuffer.FileInfo.Off + int64(wlen)

}
