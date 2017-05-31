package main

import (
	"log"

	"github.com/boltdb/bolt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type MMt struct {
	Name string
	Age  int
}

func main() {
	db, err := bolt.Open("../test_bolt/mdb", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("test"))
		v := bucket.Get([]byte("masato-100"))
		var uname MMt
		msgpack.Unmarshal(v, &uname)
		log.Printf("aaaa -> %v", uname)
		return nil
	})
	defer db.Close()
}
