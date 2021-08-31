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
	baseFSTime, baseFCCTime time.Time
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

	fs := &model.FlightSession{
		DroneID:    fl.DroneID,
		Hash:       fl.SessionToken,
		LastUpdate: baseFSTime,
	}

	res, err := db.Model(fs).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Println("FlightSession inserted")
	}

	Merge(geos, fccs, &rtts)
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
				currFCCTime := timemgr.UnixTime(job.FCCTime)
				lastUpdate := baseFSTime.Add(currFCCTime.Sub(baseFCCTime))

				fsr := &model.FlightSessionReading{
					DroneID:         1,
					FlightSessionID: 1,
					Latitude:        job.Lat,
					Longitude:       job.Long,
					Altitude:        job.Alt,
					Roll:            job.Roll,
					Pitch:           job.Pitch,
					Yaw:             job.Yaw,
					BatVoltage:      job.BatVoltage,
					BatCurrent:      job.BatCurrent,
					BatPercent:      job.BatPercent,
					BatTemperature:  job.BatTemperature,
					Temperature:     job.Temperature,
					LastUpdate:      lastUpdate,
				}
				_, err := db.Model(fsr).OnConflict("DO NOTHING").Insert()
				if err != nil {
					panic(err)
				}
				//if res.RowsAffected() > 0 {
				//	fmt.Println("FlightSession inserted")
				//}

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

			jobs <- *rtt
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

	for i, rec := range *rtts {
		if i < 5 {
			fmt.Printf("%#v\n", rec)
		}
	}

	fmt.Println("Count Concurrent ", len(*rtts))
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

	return -1, nil
}