// Package unzip
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package unzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract will decompress a zip archive, moving all dir-files and folders
// within the zip file source to an output directory destination.
func Extract(src string, dest string) (filenames []string, err error) {
	var r *zip.ReadCloser

	r, err = zip.OpenReader(src)
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

func closeZipReader() func(r *zip.ReadCloser) {
	return func(r *zip.ReadCloser) {
		fmt.Println("Defer: closing file.")
		_ = r.Close()
	}
}
