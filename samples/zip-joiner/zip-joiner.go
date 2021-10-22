// Package main
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-24
package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	chunkedFileName   = "./01.datamgr.zip"
	tmpDir            = "./tmp"
	mergedFilePostfix = "-merged.zip"
)

const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement

func main() {
	files, err := Merge(chunkedFileName, tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", files)
}

func Merge(src, dest string) ([]string, error) {
	var parts []string
	var rootPath string
	var err error

	rootPath, sf := filepath.Split(src)
	sfExt := filepath.Ext(sf)
	sfName := strings.TrimSuffix(sf, sfExt)

	if len(rootPath) == 0 {
		rootPath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	parent := filepath.Dir(rootPath)

	fmt.Println("splitFile-dir-path:", rootPath)
	fmt.Println("splitFile-name:", sf)
	fmt.Println("parent-path:", parent)

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || path == rootPath {
			r, err := regexp.MatchString(sfName, d.Name())
			if err == nil && r {
				parts = append(parts, d.Name())
				fmt.Println("f.name:", d.Name())
			}
		} else {
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	mergedFileName := sfName + mergedFilePostfix
	mergedFilePath := filepath.Join(tmpDir, mergedFileName)
	_, err = os.Create(mergedFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// set the mergedFile to APPEND MODE!!
	// open files r and w
	mergedFile, err := os.OpenFile(mergedFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	// IMPORTANT! do not defer a mergedFile.Close when opening a mergedFile for APPEND mode!
	// defer mergedFile.Close()

	// Just information on which part of the new mergedFile we are appending
	var writePosition int64 = 0
	for i, part := range parts {
		partFile, err := os.Open(part)
		if err != nil {
			log.Fatal(err)
		}
		defer partFile.Close()

		partInfo, err := partFile.Stat()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Processing file:", partInfo.Name())

		header, err := zip.FileInfoHeader(partInfo)
		if err != nil {
			return nil, err
		}
		spew.Dump(header)

		// calculate the bytes size of each chunk
		// we are not going to rely on previous datamgr and constant
		partSize := partInfo.Size()
		partBytes := make([]byte, partSize)

		//fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + partSize

		// read into partBytes
		reader := bufio.NewReader(partFile)
		_, err = reader.Read(partBytes)
		if err != nil {
			log.Fatal(err)
		}

		// DON't USE ioutil.WriteFile, it will overwrite the previous bytes!
		// Instead, write/save buffer to disk
		// ioutil.WriteFile(mergedFileName, partBytes, os.ModeAppend)
		n, err := mergedFile.Write(partBytes)
		if err != nil {
			log.Fatal(err)
		}

		_ = mergedFile.Sync() //flush to disk

		// Free up the buffer for next cycle should not be a problem if the
		// part size is small, but can be resource hogging if the part size is huge.
		// Also, it is a good practice to clean up your own plate after eating
		partBytes = nil // reset or empty our buffer

		fmt.Println("Written ", n, " bytes")
		fmt.Println("Recombining part [", i, "] into : ", mergedFileName)
	}

	// Now, close the mergedFile
	mergedFile.Close()

	return parts, nil
}

func tmp() {

	//file, err := os.Open(src)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	//fileInfo, _ := file.Stat()
	//
	//var fileSize int64 = fileInfo.Size()
	//
	//// calculate total number of parts the file will be chunked into
	//
	//totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	//
	//fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	//
	//for i := uint64(0); i < totalPartsNum; i++ {
	//	partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
	//	partBuffer := make([]byte, partSize)
	//
	//	file.Read(partBuffer)
	//
	//	// write to disk
	//	fileName := "bigfile_" + strconv.FormatUint(i, 10)
	//	_, err := os.Create(fileName)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// write/save buffer to disk
	//	ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
	//
	//	fmt.Println("Split to : ", fileName)
	//}
	//
	//// just for fun, let's recombine back the chunked files in a new file
	//
	//newFileName := "NEWbigfile.zip"
	//_, err = os.Create(newFileName)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	////set the newFileName file to APPEND MODE!!
	//// open files r and w
	//
	//file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	//// defer file.Close()
	//
	//// just information on which part of the new file we are appending
	//var writePosition int64 = 0
	//
	//for j := uint64(0); j < totalPartsNum; j++ {
	//
	//	//read a chunk
	//	currentChunkFileName := "bigfile_" + strconv.FormatUint(j, 10)
	//
	//	newFileChunk, err := os.Open(currentChunkFileName)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	defer newFileChunk.Close()
	//
	//	chunkInfo, err := newFileChunk.Stat()
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// calculate the bytes size of each chunk
	//	// we are not going to rely on previous datamgr and constant
	//
	//	var chunkSize int64 = chunkInfo.Size()
	//	chunkBufferBytes := make([]byte, chunkSize)
	//
	//	fmt.Println("Appending at position : [", writePosition, "] bytes")
	//	writePosition = writePosition + chunkSize
	//
	//	// read into chunkBufferBytes
	//	reader := bufio.NewReader(newFileChunk)
	//	_, err = reader.Read(chunkBufferBytes)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
	//	// write/save buffer to disk
	//	//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)
	//
	//	n, err := file.Write(chunkBufferBytes)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	file.Sync() //flush to disk
	//
	//	// free up the buffer for next cycle
	//	// should not be a problem if the chunk size is small, but
	//	// can be resource hogging if the chunk size is huge.
	//	// also a good practice to clean up your own plate after eating
	//
	//	chunkBufferBytes = nil // reset or empty our buffer
	//
	//	fmt.Println("Written ", n, " bytes")
	//
	//	fmt.Println("Recombining part [", j, "] into : ", newFileName)
	//}
	//
	//// now, we close the newFileName
	//file.Close()
}
