package internal

import (
	"encoding/json"
	"errors"
	"fmt"

	"os"
)

const Blocksize = 1024

type Head struct {
	Name string
	Size int64
}
type File struct {
	Handle *os.File //文件
	Header Head
	Start  byte
	Off    int64 //文件位置偏移
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

func (b *Buffer) GetFileInfo() {

	fileInfo, err := b.Fileinfo.Handle.Stat()
	if err != nil {
		panic(err)
	}

	b.Fileinfo.Header.Name = fileInfo.Name()
	b.Fileinfo.Header.Size = fileInfo.Size()
}

func (b *Buffer) ReadBlock() ([]byte, error) {
	var bs []byte
	var err error

	if b.Fileinfo.Start == 0xAA {
		bs, err = json.Marshal(b.Fileinfo.Header)
		if err != nil {
			fmt.Println("marshal error", err)
			panic(err)
		}
		return bs, err
	}
	bsize := Blocksize
	if b.Fileinfo.Header.Size <= 0 {
		return nil, errors.New("文件大小为0")
	} else if b.Fileinfo.Header.Size < Blocksize {
		bsize = int(b.Fileinfo.Header.Size)
	} else if (b.Fileinfo.Header.Size - b.Fileinfo.Off) < int64(Blocksize) {
		bsize = int(b.Fileinfo.Header.Size - b.Fileinfo.Off)
	}

	bs = make([]byte, bsize)
	n, err := b.Fileinfo.Handle.ReadAt(bs, b.Fileinfo.Off)
	if n <= 0 {
		fmt.Println(err)
		panic(err)
		return nil, errors.New("readat error")
	}
	b.Fileinfo.Off = b.Fileinfo.Off + int64(bsize)
	return bs, err
}

func (b *Buffer) WriteBlock(bs []byte) {
	var err error

	if b.Fileinfo.Start == 0xAA {
		if err := json.Unmarshal(bs, &b.Fileinfo.Header); err != nil {
			panic(err)
		}
		b.Fileinfo.Handle, err = os.Create(b.Fileinfo.Header.Name)
		if err != nil {
			panic(err)
		}
	} else {
		if _, err := b.Fileinfo.Handle.WriteAt(bs, b.Fileinfo.Off); err != nil {
			panic(err)
		}
		b.Fileinfo.Off = b.Fileinfo.Off + int64(len(bs))
	}

}
