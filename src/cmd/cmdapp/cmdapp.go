// Package cmdapp
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package cmdapp

const (
	Name  = "go-concurrency-samples"
	Short = "Unzips a file or walks a directory and process the datamgr inside"
	Long  = `This application unzip a file containing logs and csv files that will be merged and then inserted into a database.`

	FName  = "filename"
	FShort = "f"
	FDesc  = "Zip file or directory that contains the logs"

	DName  = "destination"
	DShort = "d"
	DDesc  = "Directory where the merged file will be stored"
	DDefault = ""

	EName    = "extract"
	EShort   = "e"
	EDesc    = "Extracts logs if they have been zipped"
	EDefault = false

	MName    = "merge"
	MShort   = "m"
	MDesc    = "Merges part files if the zip file has been split"
	MDefault = false
)

const (
	VersionTemplate = "%s %s.%s\n"
	Version         = "v1.0"
	Commit          = "0"
)
