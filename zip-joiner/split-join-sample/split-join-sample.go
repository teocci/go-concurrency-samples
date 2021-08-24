// Package main
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-25
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

// Hash method returns the file hash
func Hash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

// Split method splits the files into part files of user defined lengths
func Split(filename string, splitSize int) {
	bufferSize := 1024 // 1 KB for optimal splitting
	fileStats, _ := os.Stat(filename)
	pieces := int(math.Ceil(float64(fileStats.Size()) / float64(splitSize*1048576)))
	nTimes := int(math.Ceil(float64(splitSize*1048576) / float64(bufferSize)))
	file, err := os.Open(filename)
	hashFileName := filename + "-split-hash.txt"
	hashFile, err := os.OpenFile(hashFileName, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	i := 1
	for i <= pieces {
		partFileName := filename + ".pt" + strconv.Itoa(i)
		pfile, _ := os.OpenFile(partFileName, os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Println("Creating file:", partFileName)
		buffer := make([]byte, bufferSize)
		j := 1
		for j <= nTimes {
			_, inFileErr := file.Read(buffer)
			if inFileErr == io.EOF {
				break
			}
			_, err2 := pfile.Write(buffer)
			if err2 != nil {
				log.Fatal(err2)
			}
			j++
		}
		partFileHash := Hash(partFileName)
		s := partFileName + ": " + partFileHash + "\n"
		hashFile.WriteString(s)
		pfile.Close()
		i++
	}
	s := "Original file hash: " + Hash(filename) + "\n"
	hashFile.WriteString(s)
	file.Close()
	hashFile.Close()
	fmt.Printf("Splitted successfully! Find the individual file hashes in %s", hashFileName)
}

// Join method joins the split files into one, original file
func Join(startFileName string, numberParts int) {
	a := len(startFileName)
	b := a - 4
	iFileName := startFileName[:b]
	_, err := os.Create(iFileName)
	jointFile, err := os.OpenFile(iFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	i := 1
	for i <= numberParts {
		partFileName := iFileName + ".pt" + strconv.Itoa(i)
		fmt.Println("Processing file:", partFileName)
		pfile, _ := os.Open(partFileName)
		pfileinfo, err := pfile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		pfilesize := pfileinfo.Size()
		pfileBytes := make([]byte, pfilesize)
		readSrc := bufio.NewReader(pfile)
		_, err = readSrc.Read(pfileBytes)
		if err != nil {
			log.Fatal(err)
		}
		_, err = jointFile.Write(pfileBytes)
		if err != nil {
			log.Fatal(err)
		}
		pfile.Close()
		jointFile.Sync()
		pfileBytes = nil
		i++
	}
	jointFile.Close()
	fmt.Printf("Combined successfully!")
}

func main() {
	option := flag.String("opt", "SPLIT", "The option: SPLIT / JOIN")
	filename := flag.String("file", "file.txt", "The option: SPLIT / JOIN")
	splitSize := flag.Int("size", 0, "Split size, mandatory for SPLITing files (in MB)")
	nParts := flag.Int("parts", 0, "Number of parts, mandatory for JOINing files")
	flag.Parse()

	if *option == "SPLIT" || *option == "split" {
		Split(*filename, *splitSize)
	} else if *option == "JOIN" || *option == "join" {
		Join(*filename, *nParts)
	} else {
		fmt.Println("Error! Invalid value for parameter -opt")
	}
}
