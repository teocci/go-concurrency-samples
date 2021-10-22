// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package csvmgr

import (
	"bufio"
	"bytes"
	"github.com/teocci/go-concurrency-samples/src/utfmgr"
	"io"
	"log"
	"os"
)

const lineBreak = '\n'

func LineCounter(fn string) (count int) {
	f := OpenFile(fn)
	defer CloseFile()(f)

	buf := make([]byte, bufio.MaxScanTokenSize)
	lineSep := []byte{lineBreak}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			log.Fatal(err)
		}
	}
}


func UTFBufferFile(fn string) []byte {
	f, err := utfmgr.OpenFile(fn, utfmgr.UTF8)
	if err != nil {
		log.Fatal(err)
	}
	defer utfmgr.CloseFile()(f)

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(f); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func BufferFile(fn string) []byte {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseFile()(f)

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(f); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func OpenFile(fn string) *os.File {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func CreateFile(fn string) *os.File {
	w, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}

	return w
}

func CloseFile() func(f *os.File) {
	return func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
