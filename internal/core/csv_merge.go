// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-31
package core

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/teocci/go-concurrency-samples/internal/csvmgr"
	"github.com/teocci/go-concurrency-samples/internal/data"
	"github.com/teocci/go-concurrency-samples/internal/model"
	"github.com/teocci/go-concurrency-samples/internal/timemgr"
)

var (
	baseFSTime  time.Time
	baseFCCTime time.Time

	inserts int
	total   int
)

func processCSVLogs(fl *FlightLog) {
	// open the first file
	var geos []data.GEOData
	geoBuff := csvmgr.LoadDataBuff(fl.Files[data.GEOFile])
	if err := gocsv.UnmarshalBytes(geoBuff, &geos); err != nil {
		log.Fatal(err)
	}

	// open second file
	var fccs []data.FCC
	fccBuff := csvmgr.LoadDataBuff(fl.Files[data.FCCFile])
	if err := gocsv.UnmarshalBytes(fccBuff, &fccs); err != nil {
		log.Fatal(err)
	}

	// create a file writer
	var rtts []data.RTT
	rttFN := fl.LogID + "_RTTdata"
	fmt.Println("rttFN:", rttFN)
	rttPath := filepath.Join(fl.LoggerDir, rttFN+".csv")
	w := csvmgr.CreateFile(rttPath)
	defer csvmgr.CloseFile()(w)
	fl.Files[rttFN] = rttPath
	_ = rtts

	// TODO: Generate date as 2021-08-01, 13:00:00
	baseFSTime = timemgr.GenBaseDate(fl.LogNum)
	baseFCCTime = timemgr.UnixTime(geos[0].FCCTime)
	fl.SessionToken = data.FNV64aS(baseFSTime.String())
	fmt.Println(baseFSTime.Format("2006-01-02, 15:04:05"))

	fs := &model.Flight{
		DroneID:    fl.DroneID,
		Hash:       fl.SessionToken,
		LastUpdate: baseFSTime,
	}

	fs.InsertIntoDB(db)

	Merge(geos, fccs, &rtts)
	CrunchRTTData(rtts)
}

func CrunchRTTData(rtts []data.RTT) {
	inserts = 0
	data.SortRTTByFCCTime(rtts)

	for i, r := range rtts {
		if parseNInsertIntoDB(i, r) {
			inserts++
		}
	}

	total += inserts

	fmt.Printf("CSV Recs: %d | DB Inserts: %d | Total Inserts: %d\n", len(rtts), inserts, total)
}

func Merge(geos []data.GEOData, fccs []data.FCC, rtts *[]data.RTT) {
	numWps := runtime.NumCPU()
	jobs := make(chan data.RTT, numWps)
	res := make(chan data.RTT)

	var wg sync.WaitGroup
	worker := func(jobs <-chan data.RTT, results chan<- data.RTT) {
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
		for _, geo := range geos {
			var rtt *data.RTT
			var last int
			last, rtt = findFCCData(geo, fccs, last)
			if rtt != nil {
				jobs <- *rtt
			}
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		*rtts = append(*rtts, r)
	}
}

func parseNInsertIntoDB(seq int, rtt data.RTT) bool {
	currFCCTime := timemgr.UnixTime(rtt.FCCTime)
	lastUpdate := baseFSTime.Add(currFCCTime.Sub(baseFCCTime))

	fsr := &model.FlightRecord{
		DroneID:         1,
		FlightID: 1,
		Sequence:        seq,
		Latitude:        rtt.Lat,
		Longitude:       rtt.Long,
		Altitude:        rtt.Alt,
		Roll:            rtt.Roll,
		Pitch:           rtt.Pitch,
		Yaw:             rtt.Yaw,
		BatVoltage:      rtt.BatVoltage,
		BatCurrent:      rtt.BatCurrent,
		BatPercent:      rtt.BatPercent,
		BatTemperature:  rtt.BatTemperature,
		Temperature:     rtt.Temperature,
		LastUpdate:      lastUpdate,
	}

	if seq < 5 {
		fmt.Printf("[%d]->%#v\n", seq, fsr)
	}

	return fsr.InsertIntoDB(db)
}

func findFCCData(geo data.GEOData, fccs []data.FCC, offset int) (int, *data.RTT) {
	for i := offset; i < len(fccs); i++ {
		if geo.FCCTime == fccs[i].FCCTime {
			fcc := fccs[i]

			return i, &data.RTT{
				DroneID:         1,
				FlightSessionID: 1,
				FCCTime:         geo.FCCTime,
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
		}
	}

	return 0, nil
}
