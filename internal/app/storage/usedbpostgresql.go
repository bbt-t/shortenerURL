package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bbt-t/shortenerURL/configs"

	_ "github.com/lib/pq"
)

type DBPostgres struct {
	db *sql.DB
}

func NewDBPostgres() *DBPostgres {
	/*
		Initializing the Postgresql DB.
		return: DB object
	*/
	db, err := sql.Open("postgres", configs.NewConfPG().DBUrl)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	return &DBPostgres{db}
}

func (d DBPostgres) CreateTable() {
	defer d.db.Close()
	q, _ := d.db.Prepare(`
		CREATE TABLE IF NOT EXISTS "items" (
		    "id" VARCHAR,
		    "url" VARCHAR,
		    "create_at" TIMESTAMP NOT NULL,
		    )
		`)
	if _, err := q.Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("> SCHEMA CREATED <")
}

func (d DBPostgres) SaveURL(k, v string) error {
	return nil
}

func (d DBPostgres) GetURL(k string) (string, error) {
	return "", nil
}
