// Package datamgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package datamgr

import "github.com/twotwotwo/sorts/sortutil"

// GEOByFCCTime implements sort.Interface for []GEOData based on
// the FCCTime field, for sorting GEOData records in sequence.
type GEOByFCCTime []GEOData

func (a GEOByFCCTime) Len() int      { return len(a) }
func (a GEOByFCCTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Float32Key and Float32Less make the sort handle the sign bit and sort NaN
// values to the end.  There are also Float64Key and Float64Less, and
// [Type]Key functions for int types.

// Key returns a uint64 that is lower for more southerly latitudes.
func (a GEOByFCCTime) Key(i int) uint64 {
	return sortutil.Float32Key(a[i].FCCTime)
}
func (a GEOByFCCTime) Less(i, j int) bool {
	return sortutil.Float32Less(a[i].FCCTime, a[j].FCCTime)
}