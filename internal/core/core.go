// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package core

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/dirfiles"
	"github.com/teocci/go-concurrency-samples/internal/logger"
	"github.com/teocci/go-concurrency-samples/internal/unzip"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const regexSessionNum = `([^-]*)st-logger`
const droneName = "drone-01"

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
	SessionToken string
	Files        map[string]string
}

var flightLogs map[string]*FlightLog

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

	flightLogs = map[string]*FlightLog{}

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
		d, ff := filepath.Split(path)
		p := filepath.Dir(d)
		pDir, pName := filepath.Split(p)

		var re = regexp.MustCompile(regexSessionNum)

		for i, match := range re.FindAllString(pName, -1) {
			fmt.Println(match, "found at index", i)
		}

		f := strings.Split(ff, ".")[0]
		id := strings.Split(pName, "st-logger")[0]

		println("dir:", d)
		println("file:", ff)
		println("parent:", p)
		println("parent-dir:", pDir)
		println("parent-file:", pName)
		println("id:", id)

		t, err := dirfiles.Hash(id + droneName)
		if err != nil {
			return err
		}
		token := strconv.Itoa(int(t))
		println("sessionToken:", token)

		if _, ok := flightLogs[token]; ok {
			flightLogs[token].Files[f] = path
		} else {
			var fl = new(FlightLog)
			fl.SessionToken = token
			fl.Files = map[string]string{}
			fl.Files[f] = path

			flightLogs[token] = fl
		}

		spew.Dump(flightLogs)
	}

	return nil
}




