package storage

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func createTable(db *sqlx.DB, schema string) {
	db.MustExec(schema /* SQL commands */)
	log.Println("SCHEMA CREATED")
}

func getInfo(db *sqlx.DB, k string) (string, error) {
	var result string

	err := db.Get(&result, "SELECT url FROM items WHERE id=$1", k)
	if err != nil {
		return "", err
	}
	return result, nil
}

func saveInDB(db *sqlx.DB, k, v string) error {
	info := map[string]interface{}{
		"id":        k,
		"url":       v,
		"create_at": time.Now(),
	}
	_, err := db.NamedExec(
		`INSERT INTO items (id, url, create_at) VALUES (:id, :url, :create_at)`,
		info,
	)
	if err != nil {
		log.Printf("ERRER : %v", err)
	}
	return err
}
