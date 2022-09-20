package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/zouchangfu/winVuln/pkg/detector"
	"log"
)

var globalDb *bolt.DB

func init() {
	db, err := bolt.Open("D:\\workplace-go\\local\\go-demo\\winvuln\\winVuln.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	globalDb = db
}

func GetAllCve() []*detector.WinVuln {
	var list []*detector.WinVuln
	err := globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("winVuln"))
		// Get()函数不会返回错误，如果key存在，则返回byte slice值，如果不存在就会返回nil。
		c := b.Cursor()
		var keys []string
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//fmt.Printf("key = %s,value = %s\n", k, v)
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
						fmt.Println(err)
					}
					list = append(list, winVuln)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return list
}
