// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package csvmgr

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/teocci/go-concurrency-samples/internal/data"
)

func Merge(geos []data.GEOData, fccs []data.FCC, rtts []*data.RTT) {
	numWps := runtime.NumCPU()
	jobs := make(chan data.RTT, numWps)
	res := make(chan *data.RTT)

	var wg sync.WaitGroup
	worker := func(jobs <-chan data.RTT, results chan<- *data.RTT) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}

				results <- data.ParseGEOData(job)
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
			var rtt data.RTT
			var last int
			last = findFCCData(geo, fccs, last, &rtt)

			jobs <- rtt
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

	for i, rec := range rtts {
		if i < 50 {
			fmt.Printf("%#v\n", rec)
		}
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
