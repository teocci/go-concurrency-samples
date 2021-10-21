// Package main
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-18
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/teocci/go-concurrency-samples/src/jobmgr"
)

var (
	client = &http.Client{Timeout: time.Millisecond * 15000}
	terms  = []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4}
)

type Item struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	start := time.Now()
	poolNumber := runtime.NumCPU()
	fmt.Println(poolNumber)

	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(CallApi)

	for i := range terms {
		dispatcher.Submit(jobmgr.Job{
			ID: i,
			Item: Item{
				Name:      fmt.Sprintf("JobID::%d", i),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}
	end := time.Now()
	log.Print(end.Sub(start).Seconds())
}

func CallApi(id int, job jobmgr.Job) error {
	baseURL := "https://age-of-empires-2-api.herokuapp.com/api/v1/civilization/%d"

	ur := fmt.Sprintf(baseURL, job.ID)
	req, err := http.NewRequest(http.MethodGet, ur, nil)

	if err != nil {
		//log.Printf("error creating a request for term %d :: error is %+v", num, err)
		return err
	}
	res, err := client.Do(req)

	if err != nil {
		//log.Printf("error querying for term %d :: error is %+v", num, err)
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Printf("error reading response body :: error is %+v", err)
		return err
	}

	//log.Printf("%#v", job.Item)
	log.Printf("%d  :: ok", id)
	return nil
}
