// Package data
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package data

import (
	"github.com/twotwotwo/sorts"
	"strconv"
)

type FCC struct {
	FCCTime        float32 `json:"fcc_time" csv:"FCCTime"`
	GPSTime        float32 `json:"gps_time" csv:"GPSTime"`
	Temperature    float32 `json:"temperature" csv:"Temp"`
	BatVoltage     float32 `json:"bat_voltage" csv:"Bat"`
	BatCurrent     float32 `json:"bat_current" csv:"BatCurr"`
	BatPercent     float32 `json:"bat_percent" csv:"BatPercent"`
	BatTemperature float32 `json:"bat_temperature" csv:"BatTemp"`
}

func ParseFCC(data []string) *FCC {
	fccTime, _ := strconv.ParseFloat(data[0], 64)
	gpsTime, _ := strconv.ParseFloat(data[1], 64)
	temp, _ := strconv.ParseFloat(data[2], 64)
	batVol, _ := strconv.ParseFloat(data[3], 64)
	batCurr, _ := strconv.ParseFloat(data[4], 64)
	batPct, _ := strconv.ParseFloat(data[5], 64)
	batTemp, _ := strconv.ParseFloat(data[6], 64)

	return &FCC{
		FCCTime:        float32(fccTime),
		GPSTime:        float32(gpsTime),
		Temperature:    float32(temp),
		BatVoltage:     float32(batVol),
		BatCurrent:     float32(batCurr),
		BatPercent:     float32(batPct),
		BatTemperature: float32(batTemp),
	}
}

func SortFCCByFCCTime(geos []FCC) {
	sorts.ByUint64(FCCByFCCTime(geos))
}
