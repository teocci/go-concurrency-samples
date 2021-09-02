// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package model

import "time"

type Flight struct {
	ID          int       `json:"id" csv:"id" pg:"id,pk,unique"`
	DroneID     int       `json:"drone_id" csv:"drone_id" pg:"drone_id"`
	Hash        string    `json:"hash" csv:"hash" pg:"hash,unique,notnull"`
	Mission     string    `json:"mission" csv:"mission" pg:"mission,unique,notnull"`
	MissionInfo string    `json:"mission_info" csv:"mission_info" pg:"mission_info"`
	MissionType string    `json:"mission_type" csv:"mission_type" pg:"mission_type"`
	Duration    float32   `json:"duration" csv:"duration" pg:"duration"`
	Distance    float32   `json:"distance" csv:"distance" pg:"distance"`
	Status      int       `json:"status" csv:"status" pg:"status"`
	Date        time.Time `json:"date" csv:"date" pg:"date"`
	LastUpdate  time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}

type FlightRecord struct {
	ID             int       `json:"id" csv:"id" pg:"id,pk,unique"`
	DroneID        int       `json:"drone_id" csv:"drone_id" pg:"drone_id"`
	FlightID       int       `json:"flight_id" pg:"flight_id" pg:"flight_id"`
	Sequence       int       `json:"sequence" csv:"seq" pg:"sequence"`
	Latitude       float32   `json:"latitude" csv:"lat" pg:"latitude"`
	Longitude      float32   `json:"longitude" csv:"long" pg:"longitude"`
	Altitude       float32   `json:"altitude" csv:"alt" pg:"altitude"`
	Roll           float32   `json:"roll" csv:"roll" pg:"roll"`
	Pitch          float32   `json:"pitch" csv:"pitch" pg:"pitch"`
	Yaw            float32   `json:"yaw" csv:"yaw" pg:"yaw"`
	BatVoltage     float32   `json:"battery_voltage" csv:"battery_voltage" pg:"battery_voltage"`
	BatCurrent     float32   `json:"battery_current" csv:"battery_current" pg:"battery_current"`
	BatPercent     float32   `json:"battery_percentage" csv:"battery_percentage" pg:"battery_percentage"`
	BatTemperature float32   `json:"battery_temperature" csv:"battery_temperature" pg:"battery_temperature"`
	Temperature    float32   `json:"temperature" csv:"temperature" pg:"temperature"`
	GPSStatus      float32   `json:"gps_status" csv:"gps_status" pg:"gps_status"`
	DroneStatus    float32   `json:"drone_status" csv:"drone_status" pg:"drone_status"`
	LastUpdate     time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}
