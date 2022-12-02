package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/bbt-t/shortenerURL/configs"

	_ "github.com/mattn/go-sqlite3"
)

type DBSqlite struct {
	db *sql.DB
}

func NewDBSqlite() *DBSqlite {
	/*
		Initializing the SQLite DB.
		return: DB object
	*/
	db, err := sql.Open("sqlite3", configs.NewConfSQLite().DBName)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS items (id VARCHAR(8), url VARCHAR(512), create_at TIMESTAMP NOT NULL)")
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Printf("ERRER : %s", err)
	}

	return &DBSqlite{db}
}

func (d DBSqlite) SaveURL(k, v string) error {
	/*
		Adding info to the DB.
		return: Error or nil
	*/
	stmt, _ := d.db.Prepare("INSERT INTO items (id, url, create_at) VALUES (?, ?, ?)")
	defer stmt.Close()
	_, err := stmt.Exec(k, v, time.Now())
	if err != nil {
		log.Printf("ERRER : %s", err)
	}
	return err
}

func (d DBSqlite) GetURL(k string) (string, error) {
	/*
		Search for info by ID.
		param k: id by which we search in the DB
		return: result (or "") and error (or nil)
	*/
	var result string

	stmt, _ := d.db.Prepare("SELECT url FROM items WHERE id = ?")
	defer stmt.Close()
	err := stmt.QueryRow(k).Scan(&result)
	if err != nil {
		return "", err
	}

	return result, nil
}
