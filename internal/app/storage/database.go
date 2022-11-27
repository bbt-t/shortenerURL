package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func AddSchema() {
	database, _ := sql.Open("sqlite3", "./SQLDB.db")
	defer database.Close()
	stmt, _ := database.Prepare("CREATE TABLE IF NOT EXISTS items (id VARCHAR(8), url VARCHAR(512), create_at TIMESTAMP NOT NULL)")
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Printf("ERRER : %s", err)
	}
}

func SaveNewUrlSQL(k, v string) {
	database, _ := sql.Open("sqlite3", "./SQLDB.db")
	defer database.Close()

	stmt, _ := database.Prepare("INSERT INTO items (id, url, create_at) VALUES (?, ?, ?)")
	defer stmt.Close()
	if _, err := stmt.Exec(k, v, time.Now()); err != nil {
		log.Printf("ERRER : %s", err)
	}
}

func PullOutUrlSQL(k string) (string, error) {
	var result string

	database, _ := sql.Open("sqlite3", "./SQLDB.db")
	defer database.Close()

	stmt, _ := database.Prepare("SELECT url FROM items WHERE id = ?")
	defer stmt.Close()
	err := stmt.QueryRow(k).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
