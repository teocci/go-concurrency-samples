// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-24
package core

import (
	"errors"
	"fmt"
)

const (
	errSessionIndexNotFound     = "session index not found"
	errSessionIndexNotNumerical = "session index not numerical"
	errUnableToOpenCSVFile      = "unable to open %s file: %s"
	errUnableToCreateCSVFile      = "unable to create new file: %s"
)

func ErrorSessionIndexNotFound() error {
	return errors.New(errSessionIndexNotFound)
}

func ErrorSessionIndexNotNumerical() error {
	return errors.New(errSessionIndexNotNumerical)
}

func ErrorUnableToOpenCSVFile(name, err string) error {
	return errors.New(fmt.Sprintf(errUnableToOpenCSVFile, name, err))
}

func ErrorUnableToCreateCSVFile(err string) error {
	return errors.New(fmt.Sprintf(errUnableToCreateCSVFile, err))
}
