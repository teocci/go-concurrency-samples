// Package dirfiles
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package dirfiles

import (
	"os"
)

// Exists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func Exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}