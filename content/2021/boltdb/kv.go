package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

func main() {
	tmp, err := ioutil.TempFile("", "bolt-talk-example-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	path := tmp.Name()

	// START OPEN OMIT
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// END OPEN OMIT

	// START UPDATE OMIT
	put := func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("bukkit"))
		if err != nil {
			return err
		}

		// START PUT OMIT
		if err := bucket.Put([]byte("answer"), []byte("hello")); err != nil {
			return err
		}
		// END PUT OMIT

		return nil
	}
	if err := db.Update(put); err != nil {
		log.Fatal(err)
	}
	// END UPDATE OMIT

	// START VIEW OMIT
	get := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bukkit"))

		val := bucket.Get([]byte("answer"))
		if val == nil {
			// not found
			return errors.New("no answer")
		}
		fmt.Println(val)
		fmt.Println(string(val))

		return nil
	}
	if err := db.View(get); err != nil {
		log.Fatal(err)
	}
	// END VIEW OMIT
}
