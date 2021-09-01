// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-29
package data

import (
	"fmt"
	"github.com/twotwotwo/sorts"
	"strconv"
)

type GEOData struct {
	FCCTime float32 `json:"fcc_time" csv:"FCCTime"`
	Lat     float32 `json:"lat" csv:"Lat"`
	Long    float32 `json:"long" csv:"Long"`
	Alt     float32 `json:"alt" csv:"Alt"`
	Roll    float32 `json:"roll" csv:"Roll"`
	Pitch   float32 `json:"pitch" csv:"Pitch"`
	Yaw     float32 `json:"yaw" csv:"Yaw"`
}

func ParseGEOData(data []string) *GEOData {
	fccTime, _ := strconv.ParseFloat(data[0], 64)
	lat, _ := strconv.ParseFloat(data[1], 64)
	long, _ := strconv.ParseFloat(data[2], 64)
	alt, _ := strconv.ParseFloat(data[3], 64)
	roll, _ := strconv.ParseFloat(data[4], 64)
	pitch, _ := strconv.ParseFloat(data[5], 64)
	yaw, _ := strconv.ParseFloat(data[6], 64)

	return &GEOData{
		FCCTime: float32(fccTime),
		Lat:     float32(lat),
		Long:    float32(long),
		Alt:     float32(alt),
		Roll:    float32(roll),
		Pitch:   float32(pitch),
		Yaw:     float32(yaw),
	}
}

func (c GEOData) String() string {
	return fmt.Sprintf("%#v", c)
}

func SortGEOByFCCTime(geos []GEOData) {
	sorts.ByUint64(GEOByFCCTime(geos))
}