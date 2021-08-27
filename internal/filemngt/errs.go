// Package filemngt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
package filemngt

import (
	"errors"
	"fmt"
)

const (
	errCanNotFindPWD           = "cannot find the pwd: %s -> %s"
	errCanNotExpandPath        = "path %s cannot be expanded: %s"
	errCanNotMakeDir           = "cannot make directory: %s -> %s"
	errNotValidFile            = "%s is not valid"
	errCanNotOpenFile          = "%s file cannot be opened -> %s"
	errModeNotDefined          = "%s mode not defined"
	errEmptyEnvironment        = "HOMEDRIVE, HOMEPATH, or USERPROFILE are blank"
	errEmptyOutputForHomeDir   = "blank output when reading home directory"
	errPathExistIsNotDirectory = "%s path exists but is not a directory"
	errPathIsNotDirectory      = "%s path exists but is not a directory"
)

func ErrCanNotFindPWD(p, e string) error {
	return errors.New(fmt.Sprintf(errCanNotFindPWD, p, e))
}

func ErrCanNotOpenFile(f, e string) error {
	return errors.New(fmt.Sprintf(errCanNotOpenFile, f, e))
}

func ErrCanNotMakeDir(d, e string) error {
	return errors.New(fmt.Sprintf(errCanNotMakeDir, d, e))
}

func ErrNotValidFile(f string) error {
	return errors.New(fmt.Sprintf(errNotValidFile, f))
}

func ErrModeNotDefined(m string) error {
	return errors.New(fmt.Sprintf(errModeNotDefined, m))
}

func ErrEmptyEnvironment() error {
	return errors.New(errEmptyEnvironment)
}

func ErrEmptyOutputForHomeDir() error {
	return errors.New(errEmptyOutputForHomeDir)
}

func ErrPathExistIsNotDirectory(d string) error {
	return errors.New(fmt.Sprintf(errPathExistIsNotDirectory, d))
}

func ErrPathIsNotDirectory(d string) error {
	return errors.New(fmt.Sprintf(errPathIsNotDirectory, d))
}

func ErrCanNotExpandPath(p string, e string) error {
	return errors.New(fmt.Sprintf(errCanNotExpandPath, p, e))
}
