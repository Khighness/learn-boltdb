package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// @Author KHighness
// @Update 2022-11-01

func main() {
	db, err := bolt.Open("1.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
