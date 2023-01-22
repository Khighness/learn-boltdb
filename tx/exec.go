package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// @Author KHighness
// @Update 2022-11-01

func main() {
	db, err := bolt.Open("1.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update()
}

func transaction(db *bolt.DB, fn func(*bolt.Tx) error) error {
	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
