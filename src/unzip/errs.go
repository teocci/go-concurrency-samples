// Package unzip
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
package unzip

import (
	"errors"
	"fmt"
)

const (
	errCanFindPWD = "cannot find the pwd: %s"
	errFileCannotBeOpened = "file cannot be opened: %s"
)

func ErrFileCannotBeOpened(e string) error {
	return errors.New(fmt.Sprintf(errFileCannotBeOpened, e))
}

func ErrCanFindPWD(e string) error {
	return errors.New(fmt.Sprintf(errCanFindPWD, e))
}