// Package filemngt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
//go:build windows
// +build windows

package filemngt

import "os"

func CustomHomeDir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != emptyString {
		return home, nil
	}

	// Prefer standard environment variable USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != emptyString {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == emptyString || path == emptyString {
		return emptyString, ErrEmptyEnvironment()
	}

	return home, nil
}
