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
	"github.com/teocci/go-concurrency-samples/internal/gcs"
	"github.com/teocci/go-concurrency-samples/internal/model"
	"github.com/teocci/go-concurrency-samples/internal/timemgr"
)

var (
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
	flightDate := timemgr.GenBaseDate(int(fl.LogNum))
	fl.SessionToken = data.FNV64aS(flightDate.String())
	fs := &model.Flight{
		DroneID:    fl.DroneID,
		Hash:       fl.SessionToken,
		LastUpdate: flightDate,
	}

	if fs.Insert(db) {
		fmt.Println("Flight date:", fs.Date.Format("2006-01-02, 15:04:05"))
		rtts = Merge(geos, fccs, fs)
		CrunchRTTData(rtts, fs)

		fs.Update(db)
	} else {
		log.Printf("flight session: %#v was not processed", fs.Hash)
	}
}

func CrunchRTTData(rtts []data.RTT, fs *model.Flight) {
	inserts = 0

	data.SortRTTByFCCTime(rtts)
	for seq, r := range rtts {
		var prevRTT data.RTT
		if seq == 0 {
			baseFCCTime = timemgr.UnixTime(r.FCCTime)
		}
		if seq > 0 {
			prevRTT = rtts[seq-1]
		}
		fr, ok := parseNInsertIntoDB(seq, r, prevRTT, fs)
		if ok {
			fs.Length++
			fs.Duration += fr.Duration
			fs.Distance += fr.Distance
			inserts++
		}
	}

	fmt.Printf("CSV Recs: %d | DB Inserts: %d | Total Inserts: %d\n", len(rtts), inserts, total)

	total += inserts
}

func Merge(geos []data.GEOData, fccs []data.FCC, fs *model.Flight) (rtts []data.RTT) {
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
		for i, geo := range geos {
			var rtt *data.RTT
			rtt = findFCCData(geo, fccs)
			if rtt != nil {
				rtt.DroneID = fs.DroneID
				rtt.FlightID = fs.ID
				_ = i
				//if i < 10 {
				//	fmt.Printf("rtt.FCCTime[%.2f],, geo.FCCTime[%.2f]\n", rtt.FCCTime, geo.FCCTime)
				//}

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
		rtts = append(rtts, r)
	}

	return rtts
}

func parseNInsertIntoDB(seq int, currRTT data.RTT, prevRTT data.RTT, fs *model.Flight) (model.FlightRecord, bool) {
	currFCCTime := timemgr.UnixTime(currRTT.FCCTime)
	lastUpdate := fs.Date.Add(currFCCTime.Sub(baseFCCTime))

	var prevFCCTime time.Time
	var duration int64
	var distance float32
	var speed float32

	if seq > 0 {
		prevFCCTime = timemgr.UnixTime(prevRTT.FCCTime)
		duration = int64(currFCCTime.Sub(prevFCCTime) / time.Millisecond)

		orig := gcs.SCS{Lat: float64(prevRTT.Lat), Lon: float64(prevRTT.Long)}
		dest := gcs.SCS{Lat: float64(currRTT.Lat), Lon: float64(currRTT.Long)}
		distance = float32(orig.MetersTo(dest))

		if duration > 0 {
			speed = distance / float32(duration)
		}
	}

	fsr := model.FlightRecord{
		DroneID:        1,
		FlightID:       fs.ID,
		Sequence:       int64(seq),
		Duration:       duration,
		Distance:       distance,
		Speed:          speed,
		Latitude:       currRTT.Lat,
		Longitude:      currRTT.Long,
		Altitude:       currRTT.Alt,
		Roll:           currRTT.Roll,
		Pitch:          currRTT.Pitch,
		Yaw:            currRTT.Yaw,
		BatVoltage:     currRTT.BatVoltage,
		BatCurrent:     currRTT.BatCurrent,
		BatPercent:     currRTT.BatPercent,
		BatTemperature: currRTT.BatTemperature,
		Temperature:    currRTT.Temperature,
		LastUpdate:     lastUpdate,
	}

	//if seq < 5 {
	//	fmt.Printf("[%d]->%#v\n", seq, fsr)
	//	fmt.Printf("duration[%d], distance[%.2f], speed[%.2f], currFCCTime[%d], baseFCCTime[%d]\n", duration, distance, speed, currFCCTime, baseFCCTime)
	//}

	return fsr, fsr.Insert(db)
}

func findFCCData(geo data.GEOData, fccs []data.FCC) *data.RTT {
	for i := 0; i < len(fccs); i++ {
		if geo.FCCTime == fccs[i].FCCTime {
			fcc := fccs[i]

			return &data.RTT{
				FCCTime:    fcc.FCCTime,
				Lat:        geo.Lat,
				Long:       geo.Long,
				Alt:        geo.Alt,
				Roll:       geo.Roll,
				Pitch:      geo.Pitch,
				Yaw:        geo.Yaw,
				BatVoltage: fcc.BatVoltage,
				BatCurrent: fcc.BatCurrent,
				BatPercent: fcc.BatPercent,
				BatTemperature:  fcc.BatTemperature,
				Temperature:     fcc.Temperature,
				GPSTime:         fcc.GPSTime,
			}
		}
	}

	return nil
}
