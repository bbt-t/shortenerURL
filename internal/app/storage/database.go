package storage

import (
	"database/sql"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type DBSqlite struct {
	db *sql.DB
}

func NewDBSqlite() *DBSqlite {
	db, err := sql.Open("sqlite3", configs.NewConfSQLite().DBName)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	return &DBSqlite{db}
}

func (d DBSqlite) CreateSchema() {
	stmt, _ := d.db.Prepare("CREATE TABLE IF NOT EXISTS items (id VARCHAR(8), url VARCHAR(512), create_at TIMESTAMP NOT NULL)")
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Printf("ERRER : %s", err)
	}
}

func (d DBSqlite) SaveURL(k, v string) error {
	stmt, _ := d.db.Prepare("INSERT INTO items (id, url, create_at) VALUES (?, ?, ?)")
	defer stmt.Close()
	_, err := stmt.Exec(k, v, time.Now())
	if err != nil {
		log.Printf("ERRER : %s", err)
	}
	return err
}

func (d DBSqlite) GetURL(k string) (string, error) {
	var result string

	stmt, _ := d.db.Prepare("SELECT url FROM items WHERE id = ?")
	defer stmt.Close()
	err := stmt.QueryRow(k).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
