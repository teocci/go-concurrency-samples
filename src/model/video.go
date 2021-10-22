// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-22
package model

import (
	"fmt"
	gopg "github.com/go-pg/pg/v10"
	"time"
)

type Video struct {
	ID          int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	FlightID    int64     `json:"flight_id" csv:"flight_id" pg:"flight_id,notnull"`
	Name        string    `json:"name" csv:"name" pg:"name,unique,notnull"`
	Path        string    `json:"path" csv:"path" pg:"path,unique,notnull"`
	Status      int       `json:"status" csv:"status" pg:"status"`
	LastUpdate  time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}


func (fs *Video) Insert(db *gopg.DB) bool {
	res, err := db.Model(fs).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		fmt.Printf("Video[%d] inserted.\n", fs.ID)
		return true
	} else {
		err = db.Model(fs).Where("flight_id = ?", fs.FlightID).Select()
		if err != nil {
			return false
		}

		if fs.ID > 0 {
			fmt.Printf("Video[%d] exits.\n", fs.ID)
			return true
		}
	}

	return false
}

func (fs *Video) Update(db *gopg.DB) bool {
	res, err := db.Model(fs).WherePK().Update()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Printf("Video[%d]  updated.\n", fs.ID)
		return true
	}

	return false
}
