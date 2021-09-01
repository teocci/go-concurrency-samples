// Package main
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	tildeChar = '~'
	dotChar   = '.'

	emptyString = ""
	tildeString = "~"
	dotString   = "."
)

func TestFilepathJoin(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	cases := []struct {
		Base string
		Dest string
	}{
		{
			"/var/www",
			"/foo",
		},

		{
			"./split-test-files",
			"~/foo",
		},

		{
			"~/foo",
			"./split-test-files",
		},


		{
			"~/foo/bar",
			"./tmp",
		},

		{
			emptyString,
			emptyString,
		},

		{
			dotString,
			dotString,
		},

		{
			"~foo/foo",
			"",
		},
	}

	_ = homeDir
	_ = pwd

	for _, tc := range cases {
		output := filepath.Join(tc.Base, tc.Dest)

		fmt.Printf("Base: %#v\nDest: %#v\nOutput: %#v\n", tc.Base, tc.Dest, output)
		fmt.Println("------------")

	}
}
