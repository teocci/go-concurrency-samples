// Package filemngt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
package filemngt

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const (
	tildeChar = '~'
	dotChar   = '.'

	emptyString = ""
	tildeString = "~"
	dotString   = "."
)

type ExtractMode int

const (
	NativeEMode ExtractMode = iota
	UserEMode
	CustomEMode
	UndefinedEMode
)

const (
	NativeEModeTag    = "native"
	UserEModeTag      = "user"
	CustomEModeTag    = "custom"
	UndefinedEModeTag = "undefined"
)

func StartWithTilde(s string) bool {
	return len(s) > 0 && s[0] == tildeChar
}

func StartWithDot(s string) bool {
	return len(s) > 0 && s[0] == dotChar
}

func StartWithPS(s string) bool {
	return len(s) > 0 && os.IsPathSeparator(s[0])
}

func IsTilde(s string) bool {
	return len(s) == 1 && (s == tildeString || s[0] == tildeChar)
}

func IsDot(s string) bool {
	return len(s) == 1 && (s == dotString || s[0] == dotChar)
}

func IsPathSeparator(s string) bool {
	return len(s) == 1 && os.IsPathSeparator(s[0])
}

func PWD() string {
	if d, err := os.Getwd(); err != nil {
		return emptyString
	} else {
		return d
	}
}

func DirExtractPathE(dir string) (string, error) {
	return dirExtractPathE(dir, NativeEMode)
}

// dirExtractPathE expands the path to include the home directory.
// If is the path contains an extension will return a ErrPathIsNotDirectory
// error.
// This function also return errors encounter while extracting the path.
func dirExtractPathE(dir string, mode ExtractMode) (string, error) {
	if len(dir) == 0 {
		return emptyString, nil
	}

	if StartWithTilde(dir) {
		if path, err := expandByModeE(dir, mode); err == nil {
			return path, nil
		} else {
			return emptyString, ErrCanNotExpandPath(path, err.Error())
		}
	}

	if IsDot(dir) {
		return dir, nil
	}

	ext := filepath.Ext(dir)
	if len(ext) > 0 {
		return emptyString, ErrPathIsNotDirectory(dir)
	}

	return dir, nil
}

// FileParentDirE retrieve the file's parent directory.
// Uses the NativeEMode as extraction mode.
// This function also return errors if any.
func FileParentDirE(f string) (string, error) {
	return fileParentDirByModeE(f, NativeEMode)
}

func fileParentDirByModeE(fn string, m ExtractMode) (string, error) {
	if len(fn) == 0 {
		return emptyString, nil
	}

	var bPath string
	var err error

	bPath = filepath.Dir(fn)

	if IsTilde(bPath) || IsTilde(fn) || StartWithTilde(bPath) {
		if IsTilde(fn) {
			bPath = tildeString
		}
		bPath, err = ExpandByModeE(bPath, m)
		if err != nil {
			return emptyString, ErrCanNotExpandPath(bPath, err.Error())
		}

	}

	if len(bPath) == 0 ||
		IsPathSeparator(bPath) ||
		IsDot(bPath) {
		bPath, err = os.Getwd()
		if err != nil {
			return emptyString, ErrCanNotFindPWD(bPath, err.Error())
		}
	}

	vPath := filepath.VolumeName(bPath)
	if len(vPath) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return emptyString, ErrCanNotFindPWD(pwd, err.Error())
		}
		bPath = filepath.Join(pwd, bPath)
	}

	return bPath, nil
}

// FilePathE retrieve the file's path.
// This function also return errors if any.
func FilePathE(f string) (string, error) {
	return filePathE(f, NativeEMode)
}

func filePathE(f string, m ExtractMode) (string, error) {
	var basePath string
	var fPath string
	var err error
	var needExpand = false

	if len(f) == 0 {
		return emptyString, nil
	}

	if StartWithTilde(f) {
		needExpand = true
	}

	_, ffn := filepath.Split(f)

	basePath, err = fileParentDirByModeE(f, m)
	if err != nil {
		return emptyString, err
	}

	if needExpand && ffn == tildeString {
		ffn = emptyString
	}

	fPath = filepath.Join(basePath, ffn)

	//// ffn = fnName + fnExt
	//fnExt := filepath.Ext(ffn)
	//fn := strings.TrimSuffix(ffn, fnExt)
	//fmt.Println("basePath:", basePath)
	//fmt.Println("full-filename", ffn)
	//fmt.Println("filename-extension:", fnExt)
	//fmt.Println("filename", fn)
	//fmt.Println("filePath", fPath)
	//fmt.Println("------")

	return fPath, nil
}

// Expand expands a path to include the home directory as ExpandE
// Uses the NativeEMode as extraction mode.
// This function does not return errors.
func Expand(path string) string {
	return expandByMode(path, NativeEMode)
}

// ExpandByMode expands the path to include the home directory as ExpandByModeE.
// This function does not return errors.
func ExpandByMode(path string, mode ExtractMode) (string, error) {
	return expandByModeE(path, mode)
}

func expandByMode(path string, mode ExtractMode) string {
	if len(path) == 0 {
		return emptyString
	}

	if path[0] != tildeChar {
		return path
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return emptyString
	}

	if dir, err := userHomeDirE(mode); err == nil {
		return filepath.Join(dir, path[1:])
	}

	return emptyString
}

// ExpandE expands the path to include the home directory as ExpandByModeE
// Uses the NativeEMode as extraction mode.
// This function also return errors if any.
func ExpandE(path string) (string, error) {
	return expandByModeE(path, NativeEMode)
}

// ExpandByModeE expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
// This function also return errors if any.
func ExpandByModeE(path string, mode ExtractMode) (string, error) {
	return expandByModeE(path, mode)
}

func expandByModeE(path string, mode ExtractMode) (string, error) {
	if len(path) == 0 {
		return emptyString, nil
	}

	if path[0] != tildeChar {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return emptyString, errors.New("cannot expand user-specific home dir")
	}

	dir, err := userHomeDirE(mode)
	if err != nil {
		return emptyString, err
	}

	return filepath.Join(dir, path[1:]), nil
}

// UserHomeDirE returns the current user's home directory (if they have one).
// There are three extraction modes:
// 1. NativeEMode calls os.UserHomeDir
// 2. UserEMode instantiate user.Current and gets the HomeDir value
// 3. CustomEMode calls CustomHomeDir
// This function also return errors if any.
func UserHomeDirE() (string, error) {
	return userHomeDirE(NativeEMode)
}

func userHomeDirE(mode ExtractMode) (string, error) {
	switch mode {
	case NativeEMode:
		dir, err := os.UserHomeDir()
		if err != nil {
			return emptyString, err
		}

		return dir, nil
	case UserEMode:
		u, err := user.Current()
		if err != nil {
			return emptyString, err
		}

		return u.HomeDir, nil
	case CustomEMode:
		dir, err := CustomHomeDir()
		if err != nil {
			return emptyString, err
		}
		return dir, nil
	default:
		return emptyString, ErrModeNotDefined(getModeTag(UndefinedEMode))
	}
}

// UserHomeDir returns the current user's home directory (if they have one).
// Uses the NativeEMode as extraction mode.
// This function does not return errors.
func UserHomeDir() string {
	return userHomeDir(NativeEMode)
}

func userHomeDir(m ExtractMode) string {
	switch m {
	case NativeEMode:
		if d, err := os.UserHomeDir(); err == nil {
			return d
		}
	case UserEMode:
		if u, err := user.Current(); err == nil {
			return u.HomeDir
		}
	case CustomEMode:
		if d, err := CustomHomeDir(); err == nil {
			return d
		}
	}

	return emptyString
}

func getModeTag(m ExtractMode) string {
	return getModeList()[m]
}
func getModeList() []string {
	return []string{NativeEModeTag, UserEModeTag, CustomEModeTag, UndefinedEModeTag}
}
