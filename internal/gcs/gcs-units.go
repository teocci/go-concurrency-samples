// Package gcs
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-06
package gcs

type MetricLength float64

const (
	CentiMetre   MetricLength = 1
	CM                        = CentiMetre
	Metre                     = 100 * CM
	M                         = Metre
	KiloMetre                 = 1000 * Metre
	KM                        = KiloMetre
	Mile                      = 1609.344 * Metre
	Mi                        = Mile
	NauticalMile              = 1852 * Metre
	NM                        = NauticalMile
)

const (
	earthRadiusNM = earthRadiusM / NM // radius of the earth in kilometers.
	earthRadiusMi = earthRadiusM / Mi // radius of the earth in miles.
	earthRadiusKm = earthRadiusM / KM // radius of the earth in kilometers.
	earthRadiusM  = 6378137 * M       // radius of the earth in meters.
	earthRadiusCm = 637816000 * CM    // radius of the earth in centi meters.
)

func MapML() map[string]MetricLength {
	return map[string]MetricLength{"m": Metre, "KM": KM, "Mi": Mi, "NM": NM}
}

func MLTags() []string {
	siMap := MapML()
	tags := make([]string, 0, len(siMap))
	for k := range siMap {
		tags = append(tags, k)
	}

	return tags
}
