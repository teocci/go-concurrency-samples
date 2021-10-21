// Package unzip
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-26
package unzip

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

type Hashes struct {
	Hashes []HashData `json:"hashes"`
}

type HashData struct {
	Name string `json:"file_name"`
	Hash string `json:"file_hash"`
}

// Hash method returns the file hash
func Hash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
