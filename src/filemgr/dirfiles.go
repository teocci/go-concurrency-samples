// Package filemgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package filemgr

import (
	"errors"
	"hash/fnv"
	"io/ioutil"
	"os"
)

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(f string) bool {
	if stat, err := os.Stat(f); err == nil {
		return !stat.IsDir()
	}

	return false
}

// DirExists checks if a dir exists
func DirExists(p string) bool {
	if stat, err := os.Stat(p); err == nil {
		return stat.IsDir()
	}

	return false
}

func FileNotExist(p string) bool {
	// Check if file already exists
	f, err := os.Open(p)
	_ = f.Close()

	return errors.Is(err, os.ErrNotExist)
}

func IsPathValid(p string) bool {
	// Check if file already exists
	f, err := os.Open(p)
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// Attempt to create it
		var d []byte
		if err := ioutil.WriteFile(p, d, 0644); err == nil {
			_ = os.Remove(p) // And delete it
			return true
		}
	}
	_ = f.Close()

	return false
}

// MakeDirIfNotExist if the path does not exist call os.Mkdir
// to create a new directory
// This function also return errors if any.
func MakeDirIfNotExist(d string) error {
	var err error
	if _, err = os.Stat(d); os.IsNotExist(err) {
		if err = os.Mkdir(d, os.ModeDir); err == nil {
			return nil
		}
	}

	return ErrCanNotMakeDir(d, err.Error())
}

// MakeDir call os.Mkdir to create a new directory
// if already exits and check if the existing path is a directory
// This function also return errors if any.
func MakeDir(d string) error {
	var err error
	if err = os.Mkdir(d, os.ModeDir); err == nil {
		return nil
	}

	// Check that the existing path is a directory
	if os.IsExist(err) {
		var stat os.FileInfo
		if stat, err = os.Stat(d); err != nil {
			return err
		}
		if !stat.IsDir() {
			return ErrPathExistIsNotDirectory(d)
		}
		return nil
	}

	return err
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func Hash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
