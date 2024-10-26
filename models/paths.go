package models

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

func InitDB() error {
	var err error

	db, err = bolt.Open("paths.db", 0600, nil)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return err
		}

		err = root.Put([]byte("/apple"), []byte("http://www.apple.com"))
		if err != nil {
			return err
		}

		err = root.Put([]byte("/reddit"), []byte("http://www.reddit.com"))
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func GetPaths(path string) (string, error) {
	var result string

	err := db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		if root == nil {
			return errors.New("root bucket not found")
		}

		val := root.Get([]byte(path))
		if val == nil {
			return fmt.Errorf("url not found for path %v", path)
		}
		result = string(val)
		return nil
	})

	return result, err
}

func CloseDB() {
	db.Close()
}
