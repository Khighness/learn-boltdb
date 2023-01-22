package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// @Author KHighness
// @Update 2022-11-04

func main() {
	db, err := bolt.Open("1.db", 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	_ = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("bucket1"))
		return err
	})
}
