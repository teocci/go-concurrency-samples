// Package cmdapp
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package cmdapp

const (
	Name  = "go-concurrency-samples"
	Short = "Unzip a file and process the data inside"
	Long  = `This application unzip a file containing logs and csv files that will be merged and then inserted into a database.`

	FName  = "filename"
	FShort = "f"
	FDesc  = "Zip file that contains the logs"

	DName  = "destination"
	DShort = "d"
	DDesc  = "Directory where the merged file will be stored"
	DDefault = ""

	MName    = "merge"
	MShort   = "m"
	MDesc    = "If the zip file has been split"
	MDefault = false
)

const (
	VersionTemplate = "%s %s.%s\n"
	Version         = "v1.0"
	Commit          = "0"
)
