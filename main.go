// Package main
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package main

import (
	"github.com/teocci/go-concurrency-samples/src/core"
	"os"
)

func main() {
	s, ok := core.New(os.Args[1:])
	if !ok {
		os.Exit(1)
	}
	s.Wait()
}