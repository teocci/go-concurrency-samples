// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-31
package core

import (
	"runtime"
	"sync"

	"github.com/teocci/go-concurrency-samples/src/data"
	"github.com/teocci/go-concurrency-samples/src/jobmgr"
	"github.com/teocci/go-concurrency-samples/src/model"
)

func MergeConcurrent(geos []data.GEOData, fccs []data.FCC, fs *model.Flight) (records []data.RTT) {
	poolNumber := runtime.NumCPU()
	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(func(id int, job jobmgr.Job) error {
		//fmt.Printf("%+v\n", job.Item.(ItemRecord).Record)
		var record *data.RTT
		geo := job.Item.(data.GEOData)
		record = findFCCData(geo, fccs)
		if record != nil {
			record.DroneID = fs.DroneID
			record.FlightID = fs.ID
			//if i < 10 {
			//	fmt.Printf("record.FCCTime[%.2f],, geo.FCCTime[%.2f]\n", record.FCCTime, geo.FCCTime)
			//}

			records = append(records, *record)
		}

		return nil
	})

	for i, geo := range geos {
		dispatcher.Submit(jobmgr.Job{
			ID:   i,
			Item: geo,
		})
	}

	return records
}

func findFCCData(geo data.GEOData, fccs []data.FCC) *data.RTT {
	for _, fcc := range fccs {
		//fmt.Printf("geo.FCCTime[%.2f] - fcc.FCCTime[%.2f]\n", geo.FCCTime, fcc.FCCTime)
		if geo.FCCTime == fcc.FCCTime {
			return &data.RTT{
				FCCTime:        geo.FCCTime,
				Lat:            geo.Lat,
				Long:           geo.Long,
				Alt:            geo.Alt,
				Roll:           geo.Roll,
				Pitch:          geo.Pitch,
				Yaw:            geo.Yaw,
				BatVoltage:     fcc.BatVoltage,
				BatCurrent:     fcc.BatCurrent,
				BatPercent:     fcc.BatPercent,
				BatTemperature: fcc.BatTemperature,
				Temperature:    fcc.Temperature,
				GPSTime:        fcc.GPSTime,
			}
		}
	}

	return nil
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
