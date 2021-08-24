// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-24
package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/jszwec/csvutil"
	"github.com/teocci/go-concurrency-samples/internal/config"
)

const largeCSVFile = "../1000k.csv"

var (
	// list of channels to communicate with workers
	// Those will be accessed synchronously no mutex required
	workers = make(map[string]chan []string)

	// wg is to make sure all workers done before exiting main
	wg = sync.WaitGroup{}

	// mu used only for sequential printing, not relevant for program logic
	mu = sync.Mutex{}
)

func initCSVProcess(fl *FlightLog) {
	// wait for all workers to finish up before exit
	defer waitTilEnd()()

	// open the first file
	base, err := os.Open(fl.Files["GEOdata"])
	if err != nil {
		config.Log.Errorln(ErrorUnableToOpenCSVFile("GEOdata", err.Error()))
	}
	defer closeFile()(base)

	// open second file
	fcc, err := os.Open(fl.Files["FCC"])
	if err != nil {
		config.Log.Errorln(ErrorUnableToOpenCSVFile("FCC", err.Error()))
	}
	defer closeFile()(fcc)

	// create a file writer
	rttFile := fl.LogID + "_RTTdata"
	fmt.Println("rttFile:", rttFile)
	outFile := filepath.Join(fl.LoggerDir, rttFile+".csv")

	w, err := os.Create(outFile)
	if err != nil {
		config.Log.Errorln(ErrorUnableToCreateCSVFile(err.Error()))
	}
	defer closeFile()(w)

	fl.Files[rttFile] = outFile

	// wrap the file readers with CSV readers
	bReader := csv.NewReader(base)
	//fr := csv.NewReader(fcc)

	geoDataSlice := make([]*GEOData, 0)
	//fccSlice := make([]*FCC, 0)

	// wrap the out file writer with a CSV writer
	//cw := csv.NewWriter(w)
	//sessionDataSlice := make([]*FSessionData, 0)

	dec, err := csvutil.NewDecoder(bReader)
	if err != nil {
		log.Fatal(err)
	}

	bHeader := dec.Header()
	fmt.Println(bHeader)

	numWps := 100
	jobs := make(chan *GEOData, numWps)
	res := make(chan *GEOData)

	worker := func(jobs <-chan *GEOData, results chan<- *GEOData) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}

				results <- job
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
			defer wg.Done()
			worker(jobs, res)
		}()
	}

	go func() {
		for {
			geoData := new(GEOData)
			if err := dec.Decode(&geoData); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			jobs <- geoData

			//rec, err := bReader.Read()
			//if err == io.EOF {
			//	break
			//}
			//if err != nil {
			//	fmt.Println("ERROR: ", err.Error())
			//	break
			//}
			//jobs <- rec
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		geoDataSlice = append(geoDataSlice, r)
	}

	fmt.Println("Count Concurrent ", len(geoDataSlice))

	//for {
	//	rec, err := bReader.Read()
	//	if err != nil {
	//		if err == io.EOF {
	//			savePartitions()
	//			return
	//		}
	//		log.Fatal(err) // sorry for the panic
	//	}
	//	processCSV(rec, true)
	//}
}

func processCSV(rec []string, first bool) {
	l := len(rec)
	part := rec[l-1]

	if c, ok := workers[part]; ok {
		// send rec to workerClosure
		c <- rec
	} else {
		// if no workerClosure for the partition

		// make a chan
		nc := make(chan []string)
		workers[part] = nc

		// start workerClosure with this chan
		go workerClosure(nc, first)

		// send rec to workerClosure via chan
		nc <- rec
	}
}

func workerClosure(c chan []string, first bool) {
	// wg.Done signals to main workerClosure completion
	wg.Add(1)
	defer wg.Done()

	var part [][]string
	for {
		// wait for a rec or close(chan)
		rec, ok := <-c
		if ok {
			// save the rec
			// instead of accumulation in memory
			// this can be saved to file directly
			part = append(part, rec)
		} else {
			// channel closed on EOF

			// dump partition
			// locks ensures sequential printing
			// not a required for independent files
			mu.Lock()
			for _, p := range part {
				if first {
					fmt.Printf("%+v\n", p)
				}
			}
			mu.Unlock()

			return
		}
	}
}

// simply signals to workers to stop
func savePartitions() {
	for _, c := range workers {
		// signal to all workers to exit
		close(c)
	}
}

func waitTilEnd() func() {
	return func() {
		wg.Wait()
		fmt.Println("File processed.")
	}
}

func closeFile() func(f *os.File) {
	return func(f *os.File) {
		fmt.Println("Defer: closing file.")
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
