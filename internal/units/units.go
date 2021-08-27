// Package units
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package units

type Base2Bytes int64

type MetricBytes int64

const (
	Kibibyte Base2Bytes = 1024
	KiB                 = Kibibyte
	Mebibyte            = Kibibyte * 1024
	MiB                 = Mebibyte
	Gibibyte            = Mebibyte * 1024
	GiB                 = Gibibyte
	Tebibyte            = Gibibyte * 1024
	TiB                 = Tebibyte
	Pebibyte            = Tebibyte * 1024
	PiB                 = Pebibyte
	Exbibyte            = Pebibyte * 1024
	EiB                 = Exbibyte
)

const (
	Kilobyte MetricBytes = 1000
	KB                   = Kilobyte
	Megabyte             = Kilobyte * 1000
	MB                   = Megabyte
	Gigabyte             = Megabyte * 1000
	GB                   = Gigabyte
	Terabyte             = Gigabyte * 1000
	TB                   = Terabyte
	Petabyte             = Terabyte * 1000
	PB                   = Petabyte
	Exabyte              = Petabyte * 1000
	EB                   = Exabyte
)