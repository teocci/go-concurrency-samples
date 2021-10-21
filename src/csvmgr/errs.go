// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-24
package csvmgr

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	errInvalidType = "invalidType: %s"
)

func ErrInvalidType(t reflect.Type) error {
	return errors.New(fmt.Sprintf(errInvalidType, t))
}
