// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package model

import (
	gopg "github.com/go-pg/pg/v10"
)

func (fsr *FlightRecord) Insert(db *gopg.DB) bool {
	res, err := db.Model(fsr).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		//fmt.Println("FlightRecord inserted")
		return true
	}

	return false
}