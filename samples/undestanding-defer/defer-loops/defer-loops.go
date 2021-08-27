// Package defere_loops
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-25
package main

import (
	"fmt"
	"time"
)

const (
	hugeSize    = 8192
	mediumSize1 = 1024 - 3*128 // no split required - fast!
	mediumSize2 = 1024 - 2*128 // split between medium and small - slow!
	mediumSize3 = 1024 - 1*128 // split between huge and medium - fast!
	smallSize   = 128
)

// big frame, forces start of stack
func huge1(i int) byte {
	var x [hugeSize]byte
	x[i] = medium1(i)
	return x[2*i]
}

// medium frame, uses up most of StackExtra
func medium1(i int) byte {
	var x [mediumSize1]byte
	for k := 0; k < 100000000; k++ {
		x[i] = small(i)
	}
	return x[2*i]
}

// small frame, overflows stack and forces allocation of new one
func small(i int) byte {
	var x [smallSize]byte
	x[i] = byte(i)
	return x[2*i]
}

// same as above, slightly different medium size
func huge2(i int) byte {
	var x [hugeSize]byte
	x[i] = medium2(i)
	return x[2*i]
}
func medium2(i int) byte {
	var x [mediumSize2]byte
	for k := 0; k < 100000000; k++ {
		x[i] = small(i)
	}
	return x[2*i]
}

func huge3(i int) byte {
	var x [hugeSize]byte
	x[i] = medium3(i)
	return x[2*i]
}

func medium3(i int) byte {
	var x [mediumSize3]byte
	for k := 0; k < 100000000; k++ {
		x[i] = small(i)
	}
	return x[2*i]
}

//

// big frame, forces start of stack
func huge1Defer(i int) byte {
	var x [hugeSize]byte
	x[i] = medium1Defer(i)
	return x[2*i]
}

// medium frame, uses up most of StackExtra
func medium1Defer(i int) byte {
	var x [mediumSize1]byte
	for k := 0; k < 100000000; k++ {
		x[i] = smallDefer(i)
		defer func(i int) {
			x[i] = x[i] + 1
		}(i)
	}
	return x[2*i]
}

// small frame, overflows stack and forces allocation of new one
func smallDefer(i int) byte {
	var x [smallSize]byte
	x[i] = byte(i)
	return x[2*i]
}

// same as above, slightly different medium size
func huge2Defer(i int) byte {
	var x [hugeSize]byte
	x[i] = medium2Defer(i)
	return x[2*i]
}
func medium2Defer(i int) byte {
	var x [mediumSize2]byte
	for k := 0; k < 100000000; k++ {
		x[i] = smallDefer(i)
		defer func(i int) {
			x[i] = x[i] + 1
		}(i)
	}
	return x[2*i]
}

func huge3Defer(i int) byte {
	var x [hugeSize]byte
	x[i] = medium3Defer(i)
	return x[2*i]
}

func medium3Defer(i int) byte {
	var x [mediumSize3]byte
	for k := 0; k < 100000000; k++ {
		x[i] = smallDefer(i)
		defer func(i int) {
			x[i] = x[i] + 1
		}(i)
	}
	return x[2*i]
}

//

func main() {
	fmt.Println("WITHOUT DEFER:")
	withoutDefer()
	fmt.Println("\nWITH DEFER:")
	withDefer()
}

func withoutDefer() {
	t0 := time.Now()
	huge1(0)
	t1 := time.Now()
	huge2(0)
	t2 := time.Now()
	huge3(0)
	t3 := time.Now()
	fmt.Printf("  no split: %v\n", t1.Sub(t0))
	fmt.Printf("with split: %v\n", t2.Sub(t1))
	fmt.Printf("both split: %v\n", t3.Sub(t2))
}

func withDefer() {
	t0 := time.Now()
	huge1Defer(0)
	t1 := time.Now()
	huge2Defer(0)
	t2 := time.Now()
	huge3Defer(0)
	t3 := time.Now()
	fmt.Printf("  no split: %v\n", t1.Sub(t0))
	fmt.Printf("with split: %v\n", t2.Sub(t1))
	fmt.Printf("both split: %v\n", t3.Sub(t2))
}
