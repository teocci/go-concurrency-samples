// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package core

type FlightLog struct {
	DroneID      int
	DroneName    string
	SessionToken string
	LogID        string
	LogNum       int
	SessionDir   string
	LoggerDir    string
	Files        map[string]string
}

func (fl *FlightLog) setSessionDirIfEmpty(d string) {
	if len(fl.SessionDir) == 0 {
		fl.SessionDir = d
	}
}

func (fl *FlightLog) setLoggerDirIfEmpty(d string) {
	if len(fl.LoggerDir) == 0 {
		fl.LoggerDir = d
	}
}
