// Package filemgr
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-26
//go:build !windows
// +build !windows

package filemgr

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const(
	colonString = ":"

	plan9HomeEnv  = "home"
	unixHomeEnv  = "HOME"

	darwinHomeCMD = `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`
	cdnpwdCMD     = "cd && pwd"
)

func CustomHomeDir() (string, error) {
	homeEnv := unixHomeEnv
	if runtime.GOOS == "plan9" {
		// On plan9, env vars are lowercase.
		homeEnv = plan9HomeEnv
	}

	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != emptyString {
		return home, nil
	}

	var stdout bytes.Buffer

	// If that fails, try OS specific commands
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", darwinHomeCMD)
		cmd.Stdout = &stdout
		if err := cmd.Run(); err == nil {
			result := strings.TrimSpace(stdout.String())
			if result != emptyString {
				return result, nil
			}
		}
	} else {
		cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			// If the error is ErrNotFound, we ignore it. Otherwise, return it.
			if err != exec.ErrNotFound {
				return emptyString, err
			}
		} else {
			if passwd := strings.TrimSpace(stdout.String()); passwd != emptyString {
				// username:password:uid:gid:gecos:home:shell
				passwdParts := strings.SplitN(passwd, colonString, 7)
				if len(passwdParts) > 5 {
					return passwdParts[5], nil
				}
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()
	cmd := exec.Command("sh", "-c", cdnpwdCMD)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return emptyString, err
	}

	result := strings.TrimSpace(stdout.String())
	if result == emptyString {
		return emptyString, ErrEmptyOutputForHomeDir()
	}

	return result, nil
}
