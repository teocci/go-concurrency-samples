// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-29
package data

import (
	"fmt"
	"github.com/twotwotwo/sorts"
	"strconv"
)

const (
	fieldGEOFCCTime = "FCCTime"
	fieldGEOLat     = "Lat"
	fieldGEOLong    = "Long"
	fieldGeoAlt     = "Alt"
	fieldGeoRoll    = "Roll"
	fieldGeoPitch   = "BatPercent"
	fieldGeoYaw     = "BatTemp"
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
	fccTime, _ := strconv.ParseFloat(data[0], 32)
	lat, _ := strconv.ParseFloat(data[1], 32)
	long, _ := strconv.ParseFloat(data[2], 32)
	alt, _ := strconv.ParseFloat(data[3], 32)
	roll, _ := strconv.ParseFloat(data[4], 32)
	pitch, _ := strconv.ParseFloat(data[5], 32)
	yaw, _ := strconv.ParseFloat(data[6], 32)

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

func ParseGEOFields(data map[string]string) *GEOData {
	fccTime, _ := strconv.ParseFloat(data[fieldGEOFCCTime], 32)
	lat, _ := strconv.ParseFloat(data[fieldGEOLat], 32)
	long, _ := strconv.ParseFloat(data[fieldGEOLong], 32)
	alt, _ := strconv.ParseFloat(data[fieldGeoAlt], 32)
	roll, _ := strconv.ParseFloat(data[fieldGeoRoll], 32)
	pitch, _ := strconv.ParseFloat(data[fieldGeoPitch], 32)
	yaw, _ := strconv.ParseFloat(data[fieldGeoYaw], 32)

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

func GEOFields() []string {
	return []string{
		fieldGEOFCCTime,
		fieldGEOLat,
		fieldGEOLong,
		fieldGeoAlt,
		fieldGeoRoll,
		fieldGeoPitch,
		fieldGeoYaw,
	}
}
