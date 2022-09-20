package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/zouchangfu/winVuln/pkg/detector"
	"log"
	"os"
)

var globalDb *bolt.DB

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	db, err := bolt.Open(path+"\\data\\winVuln.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	globalDb = db
}

func GetAllCve() []*detector.WinVuln {
	var list []*detector.WinVuln
	err := globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("winVuln"))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		var keys []string
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, fmt.Sprintf("%s", k))
		}
		for _, k := range keys {
			bucket := b.Bucket([]byte(k))
			if bucket != nil {
				c := bucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					winVuln := &detector.WinVuln{}
					err := json.Unmarshal(v, winVuln)
					if err != nil {
						log.Println(err)
					}
					list = append(list, winVuln)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return list
}
