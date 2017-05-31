package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type MMt struct {
	Name string
	Age  int
}

func saveTo(bucket *bolt.Bucket, mmt *MMt) {
	b, _ := msgpack.Marshal(mmt)
	err := bucket.Put([]byte(mmt.Name), b)
	if err != nil {
		log.Printf("saveTo%v", err)
	}
}
func main() {
	db, err := bolt.Open("mdb", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("test"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err.Error())
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("test"))
		for i := 0; i <= 10; i++ {
			saveTo(bucket, &MMt{Name: fmt.Sprintf("masato-%v", i+2), Age: i})
		}
		return nil
	})
	defer db.Close()
}
