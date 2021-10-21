// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-21
package core

import (
	"fmt"
	"github.com/teocci/go-concurrency-samples/src/csvmgr"
	"github.com/teocci/go-concurrency-samples/src/data"
	"github.com/teocci/go-concurrency-samples/src/model"
	"github.com/teocci/go-concurrency-samples/src/timemgr"
	"log"
	"path/filepath"
)

func processCSVLogs(fl *FlightLog) {
	fmt.Printf("fl.LogNum[%d]\n", fl.LogNum)
	// open the first file
	var geos []data.GEOData
	//geoBuff := csvmgr.LineNormalizer(fl.Files[data.GEOFile])
	geoBuff := csvmgr.LoadDataBuff(fl.Files[data.GEOFile])
	geos = csvmgr.GEOParser(geoBuff)
	//if err := gocsv.UnmarshalBytes(geoBuff, &geos); err != nil {
	//	log.Fatal("error:", err)
	//}
	//if err := csvutil.Unmarshal(geoBuff, &geos); err != nil {
	//	log.Fatal("error:", err)
	//}
	//for i, r := range geos {
	//	if i < 1000 {
	//		fmt.Printf("[%d]fcc.FCCTime[%.2f]\n", i, r.FCCTime)
	//	}
	//}

	// open second file
	var fccs []data.FCC
	//fccBuff := csvmgr.LineNormalizer(fl.Files[data.FCCFile])
	fccBuff := csvmgr.LoadDataBuff(fl.Files[data.FCCFile])
	fccs = csvmgr.FCCParser(fccBuff)
	//if err := gocsv.UnmarshalBytes(fccBuff, &fccs); err != nil {
	//	log.Fatal("error:", err)
	//}
	//if err := csvutil.Unmarshal(fccBuff, &fccs); err != nil {
	//	log.Fatal("error:", err)
	//}
	//for i, r := range fccs {
	//	if i < 1000 {
	//		fmt.Printf("[%d]fcc.FCCTime[%.2f]\n", i, r.FCCTime)
	//	}
	//}

	fmt.Printf("Geo Records: %d | FCC Records: %d\n", len(geos), len(fccs))
	// create a file writer
	var records []data.RTT
	rttFN := fl.LogID + "_RTTdata"
	fmt.Println("rttFN:", rttFN)
	rttPath := filepath.Join(fl.LoggerDir, rttFN+".csv")
	w := csvmgr.CreateFile(rttPath)
	defer csvmgr.CloseFile()(w)

	fl.Files[rttFN] = rttPath
	// TODO: Generate date as 2021-08-01, 13:00:00
	flightDate := timemgr.GenBaseDate(int(fl.LogNum))
	fl.SessionToken = data.FNV64aS(flightDate.String())
	fs := &model.Flight{
		DroneID:    fl.DroneID,
		Hash:       fl.SessionToken,
		Status:     model.FlightStatusCreated,
		Date:       flightDate,
		LastUpdate: flightDate,
	}

	if fs.Insert(db) {
		fmt.Println("Flight date:", fs.Date.Format("2006-01-02, 15:04:05"))
		records = MergeConcurrent(geos, fccs, fs)
		CrunchRTTData(records, fs)

		fs.Update(db)
	} else {
		log.Printf("flight session: %#v was not processed", fs.Hash)
	}
}
