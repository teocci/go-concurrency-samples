// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package csvmgr

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/teocci/go-concurrency-samples/src/jobmgr"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/teocci/go-concurrency-samples/src/data"
)

type ItemRecord struct {
	Record interface{}
}

const (
	regexTrimmer = `\s+`
)

func GenericParser(in []byte, out interface{}) error {
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrInvalidType(reflect.TypeOf(out))
	}

	switch val.Type().Elem().Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return ErrInvalidType(val.Type())
	}

	typ := val.Type().Elem()

	csvr := csv.NewReader(bytes.NewReader(in))
	rows, err := csvr.ReadAll()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	c := len(rows)

	slice := reflect.MakeSlice(typ, c, c)
	poolNumber := runtime.NumCPU()

	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(func(id int, job jobmgr.Job) error {
		//fmt.Printf("%+v\n", job.Item.(ItemRecord).CSVRecord)
		slice = reflect.Append(slice, reflect.New(typ.Elem()).Elem())

		return nil
	})

	for i, row := range rows {
		fmt.Println(len(row))
		record := slice.Index(i).Addr().Interface()

		dispatcher.Submit(jobmgr.Job{
			ID: i,
			Item: ItemRecord{
				Record: record,
			},
		})
	}

	val.Elem().Set(slice.Slice3(0, c, c))

	return nil
}

func Merge(geos []data.GEOData, fccs []data.FCC, rtts []*data.RTT) {
	numWps := runtime.NumCPU()
	jobs := make(chan data.RTT, numWps)
	res := make(chan *data.RTT)

	var wg sync.WaitGroup
	worker := func(jobs <-chan data.RTT, results chan<- *data.RTT) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}

				results <- &job
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
			defer wg.Done()
			worker(jobs, res)
		}()
	}

	go func() {
		for _, geo := range geos {
			var rtt data.RTT
			var last int
			last = findFCCData(geo, fccs, last, &rtt)

			jobs <- rtt
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		rtts = append(rtts, r)
	}

	fmt.Println("Count Concurrent ", len(rtts))
}

func LineNormalizer(fn string) []byte {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseFile()(f)
	fileScanner := bufio.NewScanner(f)
	poolNumber := runtime.NumCPU()

	var buffer bytes.Buffer
	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(func(id int, job jobmgr.Job) error {
		//fmt.Printf("%+v\n", job.Item.(ItemRecord).CSVRecord)
		str := job.Item.(string)
		str = hardStringTrimmer(str) + "\n"
		buffer.WriteString(str)

		return nil
	})

	seq := 0
	for fileScanner.Scan() {
		str := fileScanner.Text()
		err := fileScanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		dispatcher.Submit(jobmgr.Job{
			ID:   seq,
			Item: str,
		})
		seq++
	}

	fmt.Println("----------")

	return buffer.Bytes()
}

func normalizeRow(in []string, size int) (out []string) {
	for i, str := range in {
		str = hardStringTrimmer(str)

		if i > size && len(str) > 0 {
			break
		}

		out = append(out, str)
	}

	return out
}

func hardStringTrimmer(str string) string {
	str = strings.ReplaceAll(str, "\uFEFF", "")
	str = strings.ToValidUTF8(str, "")
	str = strings.TrimSpace(str)
	//space := regexp.MustCompile(`,$`)
	//str = space.ReplaceAllString(str, "")
	return str
}

func reduceFields(fields map[string]string, filter []string) (s []string) {
	if len(fields) == 0 || len(filter) == 0 {
		return nil
	}
	for _, key := range filter {
		s = append(s, fields[key])
	}

	return s
}

func associateFields(row []string, header []string) (m map[string]string) {
	if len(row) != len(header) {
		return nil
	}

	m = make(map[string]string)
	for i, key := range header {
		m[key] = row[i]
	}

	return m
}

func findFCCData(geo data.GEOData, fccs []data.FCC, offset int, rtt *data.RTT) int {
	for i := offset; i < len(fccs); i++ {
		if geo.FCCTime == fccs[i].FCCTime {
			fcc := fccs[i]
			rtt = &data.RTT{
				Lat:            geo.Lat,
				Long:           geo.Long,
				Alt:            geo.Alt,
				Roll:           geo.Roll,
				Pitch:          geo.Pitch,
				Yaw:            geo.Yaw,
				BatVoltage:     fcc.BatVoltage,
				BatCurrent:     fcc.BatCurrent,
				BatPercent:     fcc.BatPercent,
				BatTemperature: fcc.BatTemperature,
				Temperature:    fcc.Temperature,
				GPSTime:        fcc.GPSTime,
			}

			return i
		}
	}

	return -1
}
