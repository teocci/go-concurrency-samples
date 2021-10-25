// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package core

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	gopg "github.com/go-pg/pg/v10"
	"github.com/teocci/go-concurrency-samples/src/config"
	"github.com/teocci/go-concurrency-samples/src/datamgr"
	"github.com/teocci/go-concurrency-samples/src/logger"
	"github.com/teocci/go-concurrency-samples/src/model"
	"github.com/teocci/go-concurrency-samples/src/unzip"
)

type ExecutionMode int

// Execution Modes.
const (
	EMNormal ExecutionMode = iota
	EMExtract
	EMMerge
)

const (
	regexSessionNum = `(?P<id>^[0-9]+)(?P<postfix>st.logger$)`
	droneName       = "drone-01"
	droneID         = 1
)

var db *gopg.DB

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

var fLogs map[string]*FlightLog

func Start(f string, d string, mode ExecutionMode) error {
	var dPath string

	dPath = d
	if len(dPath) == 0 {
		if mode == EMNormal {
			dPath = filepath.Dir(f)
		}
		dPath = config.TempDir
	}

	var fp = f

	if mode == EMMerge {
		fp = mergeFiles(f, dPath)
	}

	if mode > EMNormal {
		dPath = filepath.Join(dPath, droneName)
		extractFiles(fp, dPath)
	}

	loadFlightLogs(dPath)
	//spew.Dump(fLogs)

	// Init DataBase
	db = model.Setup()
	defer model.Close()(db)

	processFlightLogs()

	return nil
}

func extractFiles(fp string, path string) {
	fmt.Println("Unzipping files:")
	files, err := unzip.Extract(fp, path)
	if err != nil {
		log.Fatal(err)
	}
	if config.Verbose {
		fmt.Println("Unzipped dirs and files:\n", strings.Join(files, "\n"))
		fmt.Println("----------")
	}
}

func mergeFiles(f string, path string) string {
	fmt.Println("Merging files:")
	fp, _, err := unzip.Merge(f, path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------")

	return fp
}

func loadFlightLogs(path string) {
	fLogs = map[string]*FlightLog{}

	fmt.Println("Loading log files:")
	err := filepath.WalkDir(path, findLogPaths)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------")
}

func processFlightLogs() {
	fmt.Println("Process log files:")
	for _, fl := range fLogs {
		processCSVLogs(fl)
	}
	fmt.Println("----------")
}

// findLogPaths checks the dir tree and load log paths for each session
func findLogPaths(path string, f fs.DirEntry, e error) error {
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

			file := strings.Split(ff, ".")[0]
			if len(id) == 0 {
				config.Log.Infoln(ErrorSessionIndexNotFound())
			}
			var num int64
			if len(id) > 0 {
				n, err := strconv.ParseInt(id, 0, 64)
				if err != nil {
					config.Log.Errorln(ErrorSessionIndexNotNumerical())
				}
				num = n
			}

			token := datamgr.FNV32aS(id + droneName)
			//println("sessionToken:", token)

			if _, ok := fLogs[token]; ok {
				fLogs[token].Files[file] = path
			} else {
				var fl = new(FlightLog)
				fl.DroneID = droneID
				fl.DroneName = droneName
				fl.SessionToken = token
				fl.LogID = id
				fl.LogNum = num
				fl.Files = map[string]string{}
				fl.Files[file] = path
				fl.setSessionDirIfEmpty(base)
				fl.setLoggerDirIfEmpty(parent)

				fLogs[token] = fl
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
