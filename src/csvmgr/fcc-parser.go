// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-21
package csvmgr

import (
	"bytes"
	"encoding/csv"
	"github.com/teocci/go-concurrency-samples/src/data"
	"github.com/teocci/go-concurrency-samples/src/jobmgr"
	"io"
	"log"
	"runtime"
)

func FCCParser(in []byte) (records []data.FCC) {
	reader := csv.NewReader(bytes.NewReader(in))
	rows, err := reader.ReadAll()
	if err == io.EOF {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	var header []string
	poolNumber := runtime.NumCPU()
	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(func(id int, job jobmgr.Job) error {
		//fmt.Printf("%+v\n", job.Item.(ItemRecord).CSVRecord)
		item := job.Item.(ItemRecord)
		row := normalizeFCCRow(item.Record.([]string))

		if len(row) >= 79 {
			fields := associateFields(row, header)
			record := data.ParseFCCFields(fields)
			records = append(records, *record)
		}

		return nil
	})

	for i, row := range rows {
		//if i < 10 {
		//	//fmt.Printf("row: %v\n",row)
		//	spew.Dump(row)
		//}

		if i == 0 {
			header = normalizeFCCRow(row)
			continue
		}

		dispatcher.Submit(jobmgr.Job{
			ID: i,
			Item: ItemRecord{
				Record: row,
			},
		})
	}

	return records
}

func normalizeFCCRow(in []string) []string {
	return normalizeRow(in, 78)
}
