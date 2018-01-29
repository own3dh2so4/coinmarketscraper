package storage

import (
	"sync"

	"coinmarketscraper/coin"

	"github.com/globalsign/mgo"
)

type mongoStorage struct {
	*mgo.DialInfo
	concurrency int
}

func (ms mongoStorage) connect() *mgo.Session {
	var err error
	conn, err := mgo.DialWithInfo(ms.DialInfo)
	if err != nil {
		panic(err)
	}
	err = conn.Ping()
	if err != nil {
		panic(err)
	}
	err = conn.Ping()
	if err != nil {
		panic(err)
	}
	return conn
}

func (ms mongoStorage) Save(coins <-chan coin.Coin) <-chan struct{} {
	session := ms.connect()
	end := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(ms.concurrency)
	for i := 0; i < ms.concurrency; i++ {
		go func(db *mgo.Collection) {
			for coin := range coins {
				db.Insert(coin)
			}
			wg.Done()
		}(session.DB("").C(""))
	}
	go func() {
		wg.Wait()
		session.Close()
		end <- struct{}{}
	}()
	return end
}
