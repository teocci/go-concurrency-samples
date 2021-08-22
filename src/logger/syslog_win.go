// Package logger
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-23
// +build windows

package logger

import (
	"fmt"
	"io"
)

func NewSyslog(prefix string) (io.WriteCloser, error) {
	return nil, fmt.Errorf("not implemented on windows")
}
