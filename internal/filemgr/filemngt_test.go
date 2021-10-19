// Package filemgr
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
//go:build windows
// +build windows

package filemgr

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		_ = os.Setenv(key, bck)
	}

	if value != "" {
		_ = os.Setenv(key, value)
	} else {
		_ = os.Unsetenv(key)
	}

	return deferFunc
}

func TestDirExtractPathE(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_ = pwd

	cases := []struct {
		Input  string
		Output string
		Err    bool
	}{
		{
			"/foo",
			"/foo",
			false,
		},

		{
			"~/foo",
			filepath.Join(homeDir, "foo"),
			false,
		},

		{
			"~/foo/bar",
			filepath.Join(homeDir, "foo/bar"),
			false,
		},

		{
			emptyString,
			emptyString,
			false,
		},

		{
			dotString,
			dotString,
			false,
		},

		{
			tildeString,
			homeDir,
			false,
		},

		{
			"~foo/foo",
			"",
			true,
		},

		{
			"01.data.zip",
			"",
			true,
		},

		{
			"D:/Projects/Go/go-split-merge/src/01.data.zip",
			"",
			true,
		},

		{
			"tmp/01.data.zip",
			"",
			true,
		},

		{
			".tmp/01.data.zip",
			"",
			true,
		},

		{
			"./tmp/foo",
			"./tmp/foo",
			false,
		},
	}

	for _, tc := range cases {
		actual, err := DirExtractPathE(tc.Input)
		if (err != nil) != tc.Err {
			t.Fatalf("Input: %#v\n\nErr: %s", tc.Input, err)
		}

		if actual != tc.Output {
			t.Fatalf("\nInput: %#v\nOutput: %#v\nExpected: %#v", tc.Input, actual, tc.Output)
		} else {
			fmt.Printf("Input: %#v\nOutput: %#v\nExpected: %#v\n", tc.Input, actual, tc.Output)
			fmt.Println("------------")
		}
	}
}

func TestFileParentDirE(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	cases := []struct {
		Input  string
		Output string
		Err    bool
	}{
		{
			"/foo",
			pwd,
			false,
		},

		{
			"~/foo",
			homeDir,
			false,
		},

		{
			"~/foo/bar",
			filepath.Join(homeDir, "foo"),
			false,
		},

		{
			emptyString,
			emptyString,
			false,
		},

		{
			dotString,
			pwd,
			false,
		},

		{
			tildeString,
			homeDir,
			false,
		},

		{
			"~foo/foo",
			"",
			true,
		},

		{
			"01.data.zip",
			pwd,
			false,
		},

		{
			"D:/Projects/Go/go-split-merge/src/01.data.zip",
			"D:\\Projects\\Go\\go-split-merge\\src",
			false,
		},

		{
			"tmp/01.data.zip",
			filepath.Join(pwd, "tmp"),
			false,
		},
	}

	for _, tc := range cases {
		actual, err := FileParentDirE(tc.Input)
		if (err != nil) != tc.Err {
			t.Fatalf("Input: %#v\n\nErr: %s", tc.Input, err)
		}

		if actual != tc.Output {
			t.Fatalf("\nInput: %#v\nOutput: %#v\nExpected: %#v", tc.Input, actual, tc.Output)
		} else {
			fmt.Printf("Input: %#v\nOutput: %#v\nExpected: %#v\n", tc.Input, actual, tc.Output)
			fmt.Println("------------")
		}
	}
}

func TestGetFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	cases := []struct {
		Input  string
		Output string
		Err    bool
	}{
		{
			"/foo",
			filepath.Join(pwd, "foo"),
			false,
		},

		{
			"~/foo",
			filepath.Join(homeDir, "foo"),
			false,
		},

		{
			"~/foo/bar",
			filepath.Join(homeDir, "foo/bar"),
			false,
		},


		{
			emptyString,
			emptyString,
			false,
		},

		{
			".",
			pwd,
			false,
		},

		{
			"./",
			pwd,
			false,
		},

		{
			"~",
			homeDir,
			false,
		},

		{
			"~foo/foo",
			"",
			true,
		},

		{
			"01.data.zip",
			filepath.Join(pwd, "01.data.zip"),
			false,
		},

		{
			"tmp/01.data.zip",
			filepath.Join(pwd, "tmp/01.data.zip"),
			false,
		},
	}

	for _, tc := range cases {
		actual, err := FilePathE(tc.Input)
		if (err != nil) != tc.Err {
			t.Fatalf("Input: %#v\n\nErr: %s", tc.Input, err)
		}

		if actual != tc.Output {
			t.Fatalf("\nInput: %#v\nOutput: %#v\nExpected: %#v", tc.Input, actual, tc.Output)
		} else {
			fmt.Printf("Input: %#v\nOutput: %#v\nExpected: %#v\n", tc.Input, actual, tc.Output)
			fmt.Println("------------")
		}
	}

	defer patchEnv("HOME", "C:/custom/path/")()
	expected := filepath.Join("C:/", "custom", "path", "foo/bar")

	input := "~/foo/bar"
	actual, err := filePathE(input, CustomEMode)

	if err != nil {
		t.Errorf("No error is expected, got: %v", err)
	} else if actual != expected {
		t.Errorf("Expected: %v; actual: %v", expected, actual)
		t.Errorf("\nInput: %#v\nOutput: %#v\nExpected: %#v", input, actual, expected)
	} else {
		fmt.Printf("Input: %#v\nOutput: %#v\nExpected: %#v\n", input, actual, expected)
	}
}
