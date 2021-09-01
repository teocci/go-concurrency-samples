// Package data
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package data

import "github.com/twotwotwo/sorts/sortutil"

// RTTByFCCTime implements sort.Interface for []RTT based on
// the FCCTime field, for sorting FCC records in sequence.
type RTTByFCCTime []RTT

func (a RTTByFCCTime) Len() int      { return len(a) }
func (a RTTByFCCTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Float32Key and Float32Less make the sort handle the sign bit and sort NaN
// values to the end.  There are also Float64Key and Float64Less, and
// [Type]Key functions for int types.

// Key returns a uint64 that is lower for more southerly latitudes.
func (a RTTByFCCTime) Key(i int) uint64 {
	return sortutil.Float32Key(a[i].FCCTime)
}
func (a RTTByFCCTime) Less(i, j int) bool {
	return sortutil.Float32Less(a[i].FCCTime, a[j].FCCTime)
}
