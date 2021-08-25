// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/dirfiles"
	"github.com/teocci/go-concurrency-samples/internal/logger"
	"github.com/teocci/go-concurrency-samples/internal/unzip"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	regexSessionNum = `(?P<id>^[0-9]+)(?P<postfix>st.logger$)`
	droneName       = "drone-01"
)

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
	LogID        string
	LogNum       int
	SessionDir   string
	LoggerDir    string
	Files        map[string]string
}

var flightLogs map[string]*FlightLog

func Start(f string, isSplit bool) error {
	var err error
	var fName = f

	if isSplit {
		fName, _, err = unzip.Merge(f, config.TempDir)
		if err != nil {
			return err
		}
	}

	dest := filepath.Join(config.TempDir, droneName)
	files, err := unzip.Extract(fName, dest)
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

	spew.Dump(flightLogs)

	for _, fl := range flightLogs {
		//processCSVFiles(fl)
		initCSVProcess(fl)
	}

	spew.Dump(flightLogs)

	return nil
}

// process check if is a logger directory and call a job to process its logs.
func process(path string, f fs.DirEntry, e error) error {
	if e != nil {
		return e
	}

	if !f.IsDir() {
		d, ff := filepath.Split(path)
		parent := filepath.Dir(d)
		base, pName := filepath.Split(parent)

		//fmt.Println("file-dir-path:", d)
		//fmt.Println("file-name:", ff)
		//fmt.Println("parent-path:", parent)
		//fmt.Println("parent-name:", pName)
		//fmt.Println("base-dir:", base)

		var re = regexp.MustCompile(regexSessionNum)
		if re.MatchString(pName) {
			id := re.FindStringSubmatch(pName)[1]
			//fmt.Printf("SOLO: %#v\n", re.FindStringSubmatch(pName))
			//fmt.Printf("%#v\n", re.SubexpNames())

			//fmt.Println("id:", id)

			f := strings.Split(ff, ".")[0]
			if len(id) == 0 {
				config.Log.Infoln(ErrorSessionIndexNotFound())
			}
			var num int
			if len(id) > 0 {
				n, err := strconv.ParseInt(id, 0, 64)
				if err != nil {
					config.Log.Errorln(ErrorSessionIndexNotNumerical())
				}
				num = int(n)
			}

			t, err := dirfiles.Hash(id + droneName)
			if err != nil {
				return err
			}

			token := strconv.Itoa(int(t))
			//println("sessionToken:", token)

			if _, ok := flightLogs[token]; ok {
				flightLogs[token].Files[f] = path
			} else {
				var fl = new(FlightLog)
				fl.SessionToken = token
				fl.LogID = id
				fl.LogNum = num
				fl.Files = map[string]string{}
				fl.Files[f] = path
				fl.setSessionDirIfEmpty(base)
				fl.setLoggerDirIfEmpty(parent)

				flightLogs[token] = fl
			}

			//fmt.Println("----------")
		}
	}

	return nil
}

func processCSVFiles(fl *FlightLog) error {
	var err error
	// open the first file
	base, err := os.Open(fl.Files["GEOdata"])
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open GEOdata file: %s", err))
	}
	defer base.Close()

	// open second file
	fcc, err := os.Open(fl.Files["FCC"])
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open FCC file: %s", err))
	}
	defer fcc.Close()

	// create a file writer
	rttFile := fl.LogID + "_RTTdata"
	fmt.Println("rttFile:", rttFile)
	outFile := filepath.Join(fl.LoggerDir, rttFile+".csv")

	w, err := os.Create(outFile)
	if err != nil {
		log.Panic("\nUnable to create new file: ", err)
	}
	defer w.Close()

	fl.Files[rttFile] = outFile

	//// wrap the file readers with CSV readers
	//br := csv.NewReader(base)
	//fr := csv.NewReader(fcc)
	//
	//// wrap the out file writer with a CSV writer
	//cw := csv.NewWriter(w)
	//
	//// initialize the lines
	//bLine, b := readline(br)
	//if !b {
	//	log.Panic("\nNo CSV lines in file 1.")
	//}
	//fLine, b := readline(fr)
	//if !b {
	//	log.Panic("\nNo CSV lines in file 2.")
	//}
	//
	//// copy the files according to similar rules of the merge step in Mergesort
	//for {
	//	if compare(bLine, fLine) {
	//		writeline(bLine)
	//		if bLine, b = readline(br); !b {
	//			copy(fr, w)
	//			break
	//		}
	//	} else {
	//		writeline(fLine)
	//		if fLine, b = readline(fr); !b {
	//			copy(br, w)
	//			break
	//		}
	//	}
	//}
	//mergeFiles()
	return nil
}

//func readline(r *csv.Reader) ([]string, bool) {
//	line, e := r.Read()
//	if e != nil {
//		if e == io.EOF {
//			return nil, false
//		}
//		log.Panic("\nError reading file: ", e)
//	}
//
//	return line, true
//}
//
//func writeline(w csv.Writer, line []string) {
//	e := w.Write(line)
//	if e != nil {
//		log.Panic("\nError writing file: ", e)
//	}
//}
//
//func copy(r *csv.Reader, w csv.Writer) {
//	for line, b := readline(r); !b; r, b = readline(r) {
//		writeline(w, line)
//	}
//}
//
//func compare(base, fcc string) bool {
//	/* here, determine if line1 and line2 are in the correct order (line1 first)
//	   if so, return true, otherwise false
//	*/
//	return true
//}

func (fl *FlightLog) setSessionDirIfEmpty(d string) {
	if len(fl.SessionDir) == 0 {
		fl.SessionDir = d
	}
}

func (fl *FlightLog) setLoggerDirIfEmpty(d string) {
	if len(fl.LoggerDir) == 0 {
		fl.LoggerDir = d
	}
}
