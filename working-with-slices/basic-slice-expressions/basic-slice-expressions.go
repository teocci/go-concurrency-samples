// Package basic_slice_expressions
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-25
package main

import (
	"fmt"
)

func main() {
	a := [16]int{
		0, 1, 2, 3, 4,
		5, 6, 7, 8, 9,
		10, 11, 12, 13, 14,
		15,
	}

	b := [16]int{
		0, 1, 2, 3, 4,
		5, 6, 7, 8, 9,
	}

	fmt.Printf("a->  %v\n", a[2:]) // same as a[2 : len(a)]
	fmt.Printf("a->  %v\n", a[:3]) // same as a[0 : 3]
	fmt.Printf("a->  %v\n", a[:])  // same as a[0 : len(a)]
	fmt.Printf("a->  %v\n", a[1:4])
	fmt.Printf("a->  %v\n", a[:cap(a)])
	fmt.Printf("a->  %v\n", a[1:3:5])

	fmt.Println("------")

	fmt.Printf("b->  %v\n", b[2:]) // same as a[2 : len(a)]
	fmt.Printf("b->  %v\n", b[:3]) // same as a[0 : 3]
	fmt.Printf("b->  %v\n", b[:])  // same as a[0 : len(a)]
	fmt.Printf("b->  %v\n", b[1:4])
	fmt.Printf("b->  %v\n", b[:cap(a)])
	fmt.Printf("b->  %v\n", b[:len(a)])
	fmt.Printf("b->  %v\n", b[1:3:5])

	fmt.Println("------")

	letters := []string{"a", "b", "c", "d"}
	fmt.Printf("letters->  %v\n", letters[:])

	fmt.Println("------")

	var s []byte
	s = make([]byte, 5, 5)
	// s == []byte{0, 0, 0, 0, 0}
	fmt.Printf("s->  %v\n", s[:])

	fmt.Println("------")

	g := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	//g[:2] == []byte{'g', 'o'}
	//g[2:] == []byte{'l', 'a', 'n', 'g'}
	//g[:] == g
	fmt.Printf("g->  %v\n", g[:2])
	fmt.Printf("g->  %v\n", g[2:])
	fmt.Printf("g->  %v\n", g)


	//var n int
	//var err error
	//for i := 0; i < 32; i++ {
	//	nbytes, e := f.Read(buf[i:i+1])  // Read one byte.
	//	n += nbytes
	//	if nbytes == 0 || e != nil {
	//		err = e
	//		break
	//	}
	//}
}
