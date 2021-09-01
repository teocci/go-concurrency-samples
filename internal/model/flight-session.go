// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package model

import (
	"fmt"
	gopg "github.com/go-pg/pg/v10"
)

func (fs *FlightSession) InsertIntoDB(db *gopg.DB) bool {
	res, err := db.Model(fs).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Println("FlightSession inserted")
		return true
	}

	return false
}
