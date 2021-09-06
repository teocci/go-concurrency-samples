// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package model

import (
	"fmt"
	gopg "github.com/go-pg/pg/v10"
)

func (fs *Flight) Insert(db *gopg.DB) bool {
	res, err := db.Model(fs).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		fmt.Printf("Flight[%d] inserted.\n", fs.ID)
		return true
	} else {
		err = db.Model(fs).Where("hash = ?", fs.Hash).Select()
		if err != nil {
			return false
		}

		if fs.ID > 0 {
			fmt.Printf("Flight[%d] exits.\n", fs.ID)
			return true
		}
	}

	return false
}

func (fs *Flight) Update(db *gopg.DB) bool {
	res, err := db.Model(fs).WherePK().Update()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Printf("Flight[%d]  updated.\n", fs.ID)
		return true
	}

	return false
}
