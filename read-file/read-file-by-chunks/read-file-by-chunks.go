// Package main
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-25
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const redFileByChunksDir = "read-file"
const jsonFileName = "large-file.json"

func main() {
	var err error
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	pwd := filepath.Join(dir, redFileByChunksDir)
	jsonFilePath := filepath.Join(pwd, jsonFileName)

	f, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatalf("Error to read [file=%v]-> %v", jsonFileName, err.Error())
	}

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(f)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))

		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	log.Println("Bytes:", nBytes, "Chunks:", nChunks)
}
