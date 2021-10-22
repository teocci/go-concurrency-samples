// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-21
package core

import (
	"fmt"
	"time"

	"github.com/teocci/go-concurrency-samples/src/datamgr"
	"github.com/teocci/go-concurrency-samples/src/gcs"
	"github.com/teocci/go-concurrency-samples/src/model"
	"github.com/teocci/go-concurrency-samples/src/timemgr"
)

var (
	baseFCCTime time.Time

	inserts int
	total   int
)

func CrunchRTTData(records []datamgr.RTT, flight *model.Flight) {
	inserts = 0

	datamgr.SortRTTByFCCTime(records)
	for seq, r := range records {
		var prevRec datamgr.RTT
		if seq == 0 {
			baseFCCTime = timemgr.UnixTime(r.FCCTime)
		}
		if seq > 0 {
			prevRec = records[seq-1]
		}
		fr, ok := parseNInsertIntoDB(seq, r, prevRec, flight)
		if ok {
			flight.Length++
			flight.Duration += fr.Duration
			flight.Distance += fr.Distance
			inserts++
		}
	}

	if inserts > 0 {
		flight.Status |= model.FlightStatusCompleted | model.FlightStatusProcessed
	}

	fmt.Printf("CSV Recs: %d | DB Inserts: %d | Total Inserts: %d\n", len(records), inserts, total)
	total += inserts
}

func parseNInsertIntoDB(seq int, currRTT datamgr.RTT, prevRTT datamgr.RTT, fs *model.Flight) (model.FlightRecord, bool) {
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

	if seq < 10 {
		fmt.Printf("[%d]->%#v\n", seq, fsr)
		//fmt.Printf("duration[%d], distance[%.2f], speed[%.2f], currFCCTime[%d], baseFCCTime[%d]\n", duration, distance, speed, currFCCTime, baseFCCTime)
	}

	return fsr, fsr.Insert(db)
}
