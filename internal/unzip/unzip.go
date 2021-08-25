// Package unzip
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package unzip

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const mergedFilePostfix = "-merged.zip"

// Extract will decompress a zip archive, moving all dir-files and folders
// within the zip file source to an output directory destination.
func Extract(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, ErrFileCannotBeOpened(err.Error())
	}
	defer closeZipReader()(r)

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		_ = outFile.Close()
		_ = rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

func Merge(src, dest string) (string, []string, error) {
	var parts []string
	var rootPath string
	var mergedFileName string

	rootPath, sf := filepath.Split(src)
	sfExt := filepath.Ext(sf)
	sfName := strings.TrimSuffix(sf, sfExt)

	parent := filepath.Dir(rootPath)

	fmt.Println("splitFile-dir-path:", rootPath)
	fmt.Println("splitFile-name:", sf)
	fmt.Println("parent-path:", parent)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
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
		return mergedFileName, nil, err
	}

	mergedFileName = sfName + mergedFilePostfix
	mergedFilePath := filepath.Join(dest, mergedFileName)
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
			return mergedFileName, nil, err
		}
		spew.Dump(header)

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant
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

	return mergedFileName, parts, nil
}

func closeZipReader() func(r *zip.ReadCloser) {
	return func(r *zip.ReadCloser) {
		fmt.Println("Defer: closing file.")
		_ = r.Close()
	}
}
