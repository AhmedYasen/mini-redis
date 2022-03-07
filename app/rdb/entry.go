package rdb

import (
	"fmt"
	"sync"
	"time"
)

type Key string

type Values struct {
	Val     string
	Timeout int
}

type Db struct {
	Mu          sync.Mutex
	Persistence map[Key]Values
}

func (d *Db) Set(key Key, val string, timeout_ms int) bool {
	d.Persistence[key] = Values{Val: val, Timeout: timeout_ms}
	return true
}

func (d *Db) Get(key Key) Values {
	return d.Persistence[key]
}

func (d *Db) Remove(key Key) {
	delete(d.Persistence, key)
}

func DbElementTimeoutHandler(db *Db, check_cycle_sec uint) {
	for {
		for key, val := range db.Persistence {
			if val.Timeout > 0 && val.Timeout < int(time.Now().Unix()*1000) {
				fmt.Println("Deleting")
				db.Mu.Lock()
				db.Remove(key)
				db.Mu.Unlock()
			}
		}

		time.Sleep(time.Duration(check_cycle_sec) * time.Second)
	}
}
