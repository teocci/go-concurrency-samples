// Package filemngt
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package filemngt

import (
	"fmt"
	"testing"
)

func TestExits(t *testing.T) {
	cases := []struct {
		Input  string
		Output bool
	}{
		{
			"/foo",
			false,
		},

		{
			"~/foo",
			false,
		},

		{
			"..",
			false,
		},

		{
			"~/.android",
			false,
		},
	}

	for _, tc := range cases {
		result := FileExists(tc.Input)

		if result != tc.Output {
			t.Fatalf("\nInput: %#v\nOutput: %#v\nExpected: %#v", tc.Input, result, tc.Output)
		} else {
			fmt.Printf("Input: %#v\nOutput: %#v\nExpected: %#v\n", tc.Input, result, tc.Output)
			fmt.Println("------------")
		}
	}

}
