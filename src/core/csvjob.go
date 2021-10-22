// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-24
package core


//func initCSVProcess(fl *FlightLog) {
//	// open the first file
//	var geos []datamgr.GEOData
//	geoBuff := csvmgr.LoadDataBuff(fl.Files[datamgr.GEOFile])
//	if err := gocsv.UnmarshalBytes(geoBuff, &geos); err != nil {
//		log.Fatal(err)
//	}
//
//	// open second file
//	var fccs []datamgr.FCC
//	fccBuff := csvmgr.LoadDataBuff(fl.Files[datamgr.FCCFile])
//	if err := gocsv.UnmarshalBytes(fccBuff, &fccs); err != nil {
//		log.Fatal(err)
//	}
//
//	// create a file writer
//	var rtts []*datamgr.RTT
//	rttFN := fl.LogID + "_RTTdata"
//	fmt.Println("rttFN:", rttFN)
//	rttPath := filepath.Join(fl.LoggerDir, rttFN+".csv")
//	w := csvmgr.CreateFile(rttPath)
//	defer csvmgr.CloseFile()(w)
//	fl.Files[rttFN] = rttPath
//	_ = rtts
//
//	// init db
//	db = model.Setup()
//	defer db.Close()
//
//	// TODO: Generate date as 2021-08-01, 13:00:00
//	flightDate = timemgr.GenBaseDate(fl.LogNum)
//	baseFCCTime = timemgr.UnixTime(geos[0].FCCTime)
//	fl.SessionToken = datamgr.FNV64aS(flightDate.String())
//	fmt.Println(flightDate.Format("2006-01-02, 15:04:05"))
//
//	fs := &model.FlightSession{
//		DroneID:    fl.DroneID,
//		Hash:       fl.SessionToken,
//		LastUpdate: flightDate,
//	}
//
//	res, err := db.Model(fs).OnConflict("DO NOTHING").Insert()
//	if err != nil {
//		panic(err)
//	}
//	if res.RowsAffected() > 0 {
//		fmt.Println("FlightSession inserted")
//	}
//
//	for _, geo := range geos {
//		var last int
//		var rtt *datamgr.RTT
//		for j := last; j < len(fccs); j++ {
//			if geo.FCCTime == fccs[j].FCCTime {
//				fcc := fccs[j]
//				last = j
//				rtt = &datamgr.RTT{
//					DroneID:         1,
//					FlightID: 1,
//					Lat:             geo.Lat,
//					Long:            geo.Long,
//					Alt:             geo.Alt,
//					Roll:            geo.Roll,
//					Pitch:           geo.Pitch,
//					Yaw:             geo.Yaw,
//					BatVoltage:      fcc.BatVoltage,
//					BatCurrent:      fcc.BatCurrent,
//					BatPercent:      fcc.BatPercent,
//					BatTemperature:  fcc.BatTemperature,
//					Temperature:     fcc.Temperature,
//					GPSTime:         fcc.GPSTime,
//				}
//
//				rtts = append(rtts, rtt)
//			}
//		}
//	}
//
//	for i := 0; i < 5; i++ {
//		sec, dec := math.Modf(float64(rtts[i].GPSTime))
//		t := time.Unix(int64(sec), int64(dec*(1e3)))
//
//		fmt.Printf("FCCTime: %+v\n", t)
//	}
//	fmt.Println("Count Concurrent ", len(rtts))
//
//}

//func mergeData(geos []datamgr.GEOData, fccs []datamgr.FCC) {
//	for _, geo := range geos {
//		var last int
//		var rtt *datamgr.RTT
//		for j := last; j < len(fccs); j++ {
//			if geo.FCCTime == fccs[j].FCCTime {
//				fcc := fccs[j]
//				last = j
//				rtt = &datamgr.RTT{
//					DroneID:         1,
//					FlightID: 1,
//					Lat:             geo.Lat,
//					Long:            geo.Long,
//					Alt:             geo.Alt,
//					Roll:            geo.Roll,
//					Pitch:           geo.Pitch,
//					Yaw:             geo.Yaw,
//					BatVoltage:      fcc.BatVoltage,
//					BatCurrent:      fcc.BatCurrent,
//					BatPercent:      fcc.BatPercent,
//					BatTemperature:  fcc.BatTemperature,
//					Temperature:     fcc.Temperature,
//					GPSTime:         fcc.GPSTime,
//				}
//
//				jobs <- geoData
//
//				//rtts = append(rtts, rtt)
//				//
//				//for i := 0; i < 5; i++ {
//				//	sec, dec := math.Modf(float64(rtts[i].GPSTime))
//				//	t := time.Unix(int64(sec), int64(dec*(1e3)))
//				//
//				//	fmt.Printf("FCCTime: %+v\n", t)
//				//}
//			}
//		}
//	}
//}
//
//func processCSV(rec []string, first bool) {
//	l := len(rec)
//	part := rec[l-1]
//
//	if c, ok := workers[part]; ok {
//		// send rec to workerClosure
//		c <- rec
//	} else {
//		// if no workerClosure for the partition
//
//		// make a chan
//		nc := make(chan []string)
//		workers[part] = nc
//
//		// start workerClosure with this chan
//		go workerClosure(nc, first)
//
//		// send rec to workerClosure via chan
//		nc <- rec
//	}
//}
//
//func workerClosure(c chan []string, first bool) {
//	// wg.Done signals to main workerClosure completion
//	wg.Add(1)
//	defer wg.Done()
//
//	var part [][]string
//	for {
//		// wait for a rec or close(chan)
//		rec, ok := <-c
//		if ok {
//			// save the rec
//			// instead of accumulation in memory
//			// this can be saved to file directly
//			part = append(part, rec)
//		} else {
//			// channel closed on EOF
//
//			// dump partition
//			// locks ensures sequential printing
//			// not a required for independent files
//			mu.Lock()
//			for _, p := range part {
//				if first {
//					fmt.Printf("%+v\n", p)
//				}
//			}
//			mu.Unlock()
//
//			return
//		}
//	}
//}
//
//// simply signals to workers to stop
//func savePartitions() {
//	for _, c := range workers {
//		// signal to all workers to exit
//		close(c)
//	}
//}
//
//func waitTilEnd() func() {
//	return func() {
//		wg.Wait()
//		fmt.Println("File processed.")
//	}
//}
