// Package seek_positions_in_file
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-25
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	redFileByChunksDir = "read-file"
	textFileName       = "test.txt"
)

func main() {
	var err error
	// Get the current working directory
	dir, _ := os.Getwd()
	fmt.Println(dir)
	pwd := filepath.Join(dir, redFileByChunksDir)
	textFilePath := filepath.Join(pwd, textFileName)

	file, _ := os.Open(textFilePath)
	defer file.Close()

	// Offset is how many bytes to move
	// Offset can be positive or negative
	var offset int64 = 5

	// Whence is the point of reference for offset
	// 0 = Beginning of file
	// 1 = Current position
	// 2 = End of file
	var whence int = io.SeekStart
	var newPos int64

	newPos, err = file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Moved to five | [offset, whence]: newPos -> [%d, %d]: %d\n", offset, whence, newPos)

	// Go back 2 bytes from current position
	offset, whence = -2, io.SeekCurrent
	newPos, err = file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Moved back two bytes | [offset, whence]: newPos -> [%d, %d]: %d\n", offset, whence, newPos)

	// Find the current position by getting the
	// return value from Seek after moving 0 bytes
	offset, whence = 0, io.SeekCurrent
	currentPosition, err := file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current position | [offset, whence]: currentPosition -> [%d, %d]: %d\n", offset, whence,  currentPosition)

	// Go to beginning of file
	offset, whence = 0, io.SeekStart
	newPos, err = file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position after seeking [0,0] | [offset, whence]: newPos -> [%d, %d]: %d\n", offset, whence, newPos)
}
