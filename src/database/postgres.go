package database

import (
	"database/sql"
	"log"

    _ "github.com/lib/pq"
)

func InitDB(dbUrl string, production bool) *sql.DB {
    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        log.Fatal("Could not open db :: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Could not ping db :: ", err)
    }

    return db
}

