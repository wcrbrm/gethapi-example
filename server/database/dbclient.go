package database

import (
_	"database/sql"
	"github.com/jmoiron/sqlx"
_	"github.com/lib/pq"
	"log"
)

func EnsureConnected() {
	log.Println("[db] ensure connected")
}

// Connect - return a connection to a postgres database wi
func Connect(connectURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectURL)
	if err != nil {
		log.Fatal("ERROR OPENING DB, NOT INITIALIZING", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
