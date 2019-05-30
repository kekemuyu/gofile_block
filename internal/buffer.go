package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"os"
)

const Blocksize = 1024

type File struct {
	Name string
	Size int64
	Off  int64
	Data []byte
}

type Buffer struct {
	Filein   *os.File
	FileInfo File
}

var Defaultbuffer = &Buffer{
	Filein: nil,
}

func (b *Buffer) Create(fileName string) {
	fin, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	fileInfo, err := fin.Stat()
	if err != nil {
		panic(err)
	}
	fileObj := File{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Off:  0,
		Data: make([]byte, Blocksize),
	}

	Defaultbuffer = &Buffer{
		Filein:   fin,
		FileInfo: fileObj,
	}
}

func (b *Buffer) GetBytesbuffer(fileName string) bytes.Buffer {
	bs, err := ioutil.ReadFile(b.FileInfo.Name)

	if err != nil {
		panic(err)
	}

	b.FileInfo.Data = bs
	objBytes, err := json.Marshal(b.FileInfo)
	if err != nil {
		panic(err)
	}

	var data bytes.Buffer
	_, err = data.Write(objBytes)
	if err != nil {
		panic(err)
	}
	return data
}

func (b *Buffer) PutBytesbufferToFile(bs []byte) {
	//	var obj File
	//	n, err := json.Unmarshal(bs, &obj)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}

	//	ioutil.WriteFile(obj.Name, obj.Data, 0666)
}

func (b *Buffer) ReadBlock() (bytes.Buffer, error) {
	var outbb bytes.Buffer
	var bs []byte
	var err error

	bs = make([]byte, Blocksize)

	n, err := b.Filein.ReadAt(bs, b.FileInfo.Off)
	if (n <= 0) && (err != io.EOF) {
		fmt.Println(err)
		return bytes.Buffer{}, err
	}

	b.FileInfo = File{
		Data: bs,
	}
	outbs, err := json.Marshal(b.FileInfo)
	if err != nil {
		fmt.Println("ReadBlock marshal error", err)
		return bytes.Buffer{}, err
	}
	if _, err = outbb.Write(outbs); err != nil {
		fmt.Println("ReadBlock writebuffer error", err)
		return bytes.Buffer{}, err
	}

	b.FileInfo.Off = b.FileInfo.Off + int64(n)
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
	fmt.Println(fileinfo)
	if Defaultbuffer.Filein == nil {
		if Defaultbuffer.Filein, err = os.Create(fileinfo.Name); err != nil {
			fmt.Println("create file error", err)
			return
		}
		Defaultbuffer.Create(Defaultbuffer.FileInfo.Name)
		fmt.Println(Defaultbuffer)
	}
	Defaultbuffer.FileInfo.Data = fileinfo.Data
	bsize := Blocksize
	if len(Defaultbuffer.FileInfo.Data) < bsize {
		bsize = len(Defaultbuffer.FileInfo.Data)
	}

	var wlen int
	fmt.Println(Defaultbuffer.FileInfo.Off)
	if wlen, err = Defaultbuffer.Filein.WriteAt(Defaultbuffer.FileInfo.Data, Defaultbuffer.FileInfo.Off); err != nil {
		panic(err)
		return
	}
	Defaultbuffer.FileInfo.Off = Defaultbuffer.FileInfo.Off + int64(wlen)

}
