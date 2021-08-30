// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-24
package core

import (
	"bytes"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/data"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

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
	var geos []data.GEOData
	geoBuff := loadFileBuff(fl.Files[data.GEOFile])
	if err := gocsv.UnmarshalBytes(geoBuff, &geos); err != nil {
		log.Fatal(err)
	}

	// open second file
	var fccs []data.FCC
	fccBuff := loadFileBuff(fl.Files[data.FCCFile])
	if err := gocsv.UnmarshalBytes(fccBuff, &fccs); err != nil {
		log.Fatal(err)
	}

	// create a file writer
	var rtts []*data.RTT
	rttFN := fl.LogID + "_RTTdata"
	fmt.Println("rttFN:", rttFN)
	rttPath := filepath.Join(fl.LoggerDir, rttFN+".csv")
	w := createFile(rttPath)
	defer closeFile()(w)
	_ = rtts

	fl.Files[rttFN] = rttPath

	for _, geo := range geos {
		var last int
		var rtt *data.RTT
		for j := last; j < len(fccs); j++ {
			if geo.FCCTime == fccs[j].FCCTime {
				fcc := fccs[j]
				last = j
				rtt = &data.RTT{
					DroneID:         1,
					FlightSessionID: 1,
					Lat:             geo.Lat,
					Long:            geo.Long,
					Alt:             geo.Alt,
					Roll:            geo.Roll,
					Pitch:           geo.Pitch,
					Yaw:             geo.Yaw,
					BatVoltage:      fcc.BatVoltage,
					BatCurrent:      fcc.BatCurrent,
					BatPercent:      fcc.BatPercent,
					BatTemperature:  fcc.BatTemperature,
					Temperature:     fcc.Temperature,
					GPSTime:         fcc.GPSTime,
				}

				rtts = append(rtts, rtt)
			}
		}
	}



	for i := 0; i < 5; i++ {
		sec, dec := math.Modf(float64(rtts[i].GPSTime))
		t := time.Unix(int64(sec), int64(dec*(1e3)))

		fmt.Printf("FCCTime: %+v\n", t)
	}
	fmt.Println("Count Concurrent ", len(rtts))

}

func findFCCData(geo data.GEOData, fccs []data.FCC, offset int, rtt *data.RTT) int {
	for i := offset; i < len(fccs); i++ {
		if geo.FCCTime == fccs[i].FCCTime {
			fcc := fccs[i]
			rtt = &data.RTT{
				DroneID:         1,
				FlightSessionID: 1,
				Lat:             geo.Lat,
				Long:            geo.Long,
				Alt:             geo.Alt,
				Roll:            geo.Roll,
				Pitch:           geo.Pitch,
				Yaw:             geo.Yaw,
				BatVoltage:      fcc.BatVoltage,
				BatCurrent:      fcc.BatCurrent,
				BatPercent:      fcc.BatPercent,
				BatTemperature:  fcc.BatTemperature,
				Temperature:     fcc.Temperature,
				GPSTime:         fcc.GPSTime,
			}

			return i
		}
	}

	return -1
}

//func mergeData(geos []data.GEOData, fccs []data.FCC) {
//	for _, geo := range geos {
//		var last int
//		var rtt *data.RTT
//		for j := last; j < len(fccs); j++ {
//			if geo.FCCTime == fccs[j].FCCTime {
//				fcc := fccs[j]
//				last = j
//				rtt = &data.RTT{
//					DroneID:         1,
//					FlightSessionID: 1,
//					Lat:             geo.Lat,
//					Long:            geo.Long,
//					Alt:             geo.Alt,
//					Roll:            geo.Roll,
//					Pitch:           geo.Pitch,
//					Yaw:             geo.Yaw,
//					BatVoltage:      fcc.BatVoltage,
//					BatCurrent:      fcc.BatCurrent,
//					BatPercent:      fcc.BatPercent,
//					BatTemperature:  fcc.BatTemperature,
//					Temperature:     fcc.Temperature,
//					GPSTime:         fcc.GPSTime,
//				}
//
//				jobs <- geoData
//
//				//rtts = append(rtts, rtt)
//				//
//				//for i := 0; i < 5; i++ {
//				//	sec, dec := math.Modf(float64(rtts[i].GPSTime))
//				//	t := time.Unix(int64(sec), int64(dec*(1e3)))
//				//
//				//	fmt.Printf("FCCTime: %+v\n", t)
//				//}
//			}
//		}
//	}
//}

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

func loadFileBuff(f string) []byte {
	file, err := os.Open(f)
	if err != nil {
		config.Log.Errorln(ErrorUnableToOpenCSVFile(f, err.Error()))
	}
	defer closeFile()(file)

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(file); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func createFile(f string) *os.File {
	w, err := os.Create(f)
	if err != nil {
		config.Log.Errorln(ErrorUnableToCreateCSVFile(f, err.Error()))
	}

	return w
}
