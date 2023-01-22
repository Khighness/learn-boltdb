package main

import (
	"bytes"
	"encoding/binary"
)

// @Author KHighness
// @Update 2022-11-03

func main() {
	db, err := bolt.Open("2.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	store := lib.NewStore(db)

	_ = store.CleanupBuckets()
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
