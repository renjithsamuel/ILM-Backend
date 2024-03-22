package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // pq driver.
)

var (
	// quitRetry : exit point to exit from infinite loop
	quitRetry = make(chan bool)
)

// Connect2DB connects to the db after retryFrequency duration
func Connect2DB(connStr string, retryFrequency time.Duration) *sql.DB {
	for {
		select {
		case <-time.After(retryFrequency):
			DB, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Println("error connecting to database: ", err)
			} else {
				return DB
			}
			log.Printf("reconnecting after: %v milliseconds\n", retryFrequency)
		case <-quitRetry:
			return nil
		}
	}
}
