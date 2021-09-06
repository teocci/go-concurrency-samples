// Package haversine
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-06
package main

import (
	"fmt"
	"math"
)

type MetricLength float64

const (
	Metre        MetricLength = 1
	M                         = Metre
	KiloMetre                 = 1000 * Metre
	Km                        = KiloMetre
	Mile                      = 1609.344 * Metre
	Mi                        = Mile
	NauticalMile              = 1852 * Metre
	NM                        = NauticalMile
)

const (
	earthRadiusNM = earthRadiusM / NM // radius of the earth in kilometers.
	earthRadiusMi = earthRadiusM / Mi // radius of the earth in miles.
	earthRadiusKm = earthRadiusM / Km // radius of the earth in kilometers.
	earthRadiusM  = 6378137 * M       // radius of the earth in meters.
)

// Coord represents a geographic coordinate.
type Coord struct {
	Lat float64
	Lon float64
}

type Delta struct {
	Lat float64
	Lon float64
}

func (c Coord) MetresTo(r Coord) MetricLength {
	return Distance(c, r)
}

func (c Coord) Delta(r Coord) Delta {
	return Delta{
		Lat: c.Lat - r.Lat,
		Lon: c.Lon - r.Lon,
	}
}

func (c Coord) toRadians() Coord {
	return Coord{
		Lat: degreesToRadians(c.Lat),
		Lon: degreesToRadians(c.Lon),
	}
}

func Distance(orig, dest Coord) MetricLength {
	orig = orig.toRadians()
	dest = dest.toRadians()

	delta := orig.Delta(dest)

	a := math.Pow(math.Sin(delta.Lat/2), 2) + math.Cos(orig.Lat)*math.Cos(dest.Lat)*math.Pow(math.Sin(delta.Lon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return MetricLength(c) * earthRadiusM
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// distance calculates the distance between two points (given the latitude/longitude of those points).
// South latitudes are negative, east longitudes are positive
// orig = Latitude and Longitude of orig point (in decimal degrees)
// dest = Latitude and Longitude of dest point (in decimal degrees)
// unit = the unit you desire for results
// where:
// 'm' is in meters (default)
// 'Km' is in kilometers
// 'Mi' is in miles
// 'NM' is in nautical miles

func distance(orig Coord, dest Coord, unit ...string) float64 {
	radOrigLat := degreesToRadians(orig.Lat)
	radDestLat := degreesToRadians(dest.Lat)

	theta := orig.Lon - dest.Lon
	radTheta := degreesToRadians(theta)

	dist := math.Sin(radOrigLat)*math.Sin(radDestLat) + math.Cos(radOrigLat)*math.Cos(radDestLat)*math.Cos(radTheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "M" {
			dist = dist * 1609.344
		} else if unit[0] == "Km" {
			dist = dist * 1.609344
		} else if unit[0] == "NM" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func main() {
	fmt.Printf("%f Meters\n", distance(Coord{32.9697, -96.80322}, Coord{29.46786, -98.53506}, "M"))
	fmt.Printf("%f Miles\n", distance(Coord{32.9697, -96.80322}, Coord{29.46786, -98.53506}, "Mi"))
	fmt.Printf("%f Kilometers\n", distance(Coord{32.9697, -96.80322}, Coord{29.46786, -98.53506}, "Km"))
	fmt.Printf("%f Nautical Miles\n", distance(Coord{32.9697, -96.80322}, Coord{29.46786, -98.53506}, "NM"))
}
