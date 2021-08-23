// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package core

import (
	"context"
	"fmt"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/logger"
	"github.com/teocci/go-concurrency-samples/internal/unzip"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

const regexSessionNum = `([^-]*)st-logger`

// Core is an instance of rtsp-simple-server.
type Core struct {
	ctx       context.Context
	ctxCancel func()
	confPath  string
	confFound bool
	logger    *logger.Logger

	// out
	done chan struct{}
}

type FlightLog struct {
	SessionID   string
	FCCFile     string
	GEOdataFile string
}

var flightLogs map[string]FlightLog

func Start(f string) error {
	var err error

	files, err := unzip.Extract(f, config.TempDir)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("Unzipped files and dirs:\n" + strings.Join(files, "\n"))
		fmt.Println("----------")
	}

	flightLogs =  map[string]FlightLog{}

	err = filepath.WalkDir(config.TempDir, process)
	if err != nil {
		return err
	}

	return nil
}

// process check if is a logger directory and call a job to process its logs.
func process(path string, f fs.DirEntry, e error) error {
	if e != nil {
		return e
	}

	if !f.IsDir() {
		d, f := filepath.Split(path)
		p := filepath.Dir(d)
		pDir, pName := filepath.Split(p)

		var re = regexp.MustCompile(regexSessionNum)

		for i, match := range re.FindAllString(pName, -1) {
			fmt.Println(match, "found at index", i)
		}

		println("dir:", d)
		println("file:", f)
		println("parent:", p)
		println("parent-dir:", pDir)
		println("parent-file:", pName)



	}

	return nil
}
