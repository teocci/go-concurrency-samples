// Package unzip
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package unzip

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/teocci/go-concurrency-samples/src/filemgr"
)

const (
	tmpMergeDir       = "./merge"
	emptyString       = ""
	joinFilePrefix    = "merged-"
	regExFirstElement = "^"
)

func Merge(src, dest string) (string, []string, error) {
	var parts []string
	var basePath, destPath, mergedFPath string
	var srcFFN, mergedFN string
	var err error

	if filemgr.FileExists(src) {
		basePath, srcFFN = filepath.Split(src)
		srcExt := filepath.Ext(srcFFN)
		srcFN := strings.TrimSuffix(srcFFN, srcExt)

		if len(basePath) == 0 {
			basePath = filemgr.PWD()
		}
		fmt.Println("basePath:", basePath)
		//parent := filepath.Dir(basePath)
		//fmt.Println("parent:", parent)

		destPath, _ = filemgr.DirExtractPathE(dest)
		if len(destPath) == 0 {
			fmt.Println("Make working dir:", tmpMergeDir)
			destPath = filepath.Join(basePath, tmpMergeDir)
			if err = filemgr.MakeDirIfNotExist(destPath); err != nil {
				return emptyString, nil, err
			}
		}
		fmt.Println("destPath:", destPath)

		err = filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() || path == basePath {
				r, err := regexp.MatchString(regExFirstElement+srcFN, d.Name())
				if err == nil && r {
					parts = append(parts, d.Name())
					fmt.Println("part fn:", d.Name())
				}
			} else {
				return filepath.SkipDir
			}

			return nil
		})
		if err != nil {
			return mergedFN, nil, err
		}

		mergedFN = joinFilePrefix + srcFN
		mergedFPath = filepath.Join(destPath, mergedFN)
		_, err = os.Create(mergedFPath)
		if err != nil {
			log.Fatal(err)
		}

		// set the mergedFile to APPEND MODE!!
		// open files r and w
		mergedFile, err := os.OpenFile(mergedFPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT! do not defer a mergedFile.Close when opening a mergedFile for APPEND mode!
		// defer mergedFile.Close()

		// Just information on which part of the new mergedFile we are appending
		var writePosition int64 = 0
		for i, part := range parts {
			partPath := filepath.Join(basePath, part)
			partFile, err := os.Open(partPath)
			if err != nil {
				return emptyString, nil, err
			}

			partInfo, err := partFile.Stat()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Processing file:", partInfo.Name())

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
			// ioutil.WriteFile(mergedFN, partBytes, os.ModeAppend)
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
			fmt.Println("Recombining part [", i, "] into : ", mergedFN)

			partFile.Close()
		}

		// Now, close the mergedFile
		mergedFile.Close()
	}

	return mergedFPath, parts, nil
}
