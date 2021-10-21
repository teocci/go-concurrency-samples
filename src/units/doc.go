// Package units
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-06
//
// units provides helpful unit multipliers and functions for Go.
//
// The goal of this package is to have functionality similar to [the std time package][1].
//
//
// [1] http://golang.org/pkg/time/
//
// It allows for code like this:
// ```go
// n, err := ParseBase2Bytes("1KB")
// // n == 1024
// n = units.Mebibyte * 512
// ```
package units
