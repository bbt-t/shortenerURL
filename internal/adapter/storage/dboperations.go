package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func createTable(db *sqlx.DB, schema string) {
	db.MustExec(schema /* SQL commands */)
	log.Println("SCHEMA CREATED")
}

func getHashURL(db *sqlx.DB, k string) (string, error) {
	var result string

	err := db.Get(&result, "SELECT original_url FROM items WHERE short_url=$1", k)
	return result, err
}

func addNewUser(db *sqlx.DB, userID uuid.UUID) {
	info := map[string]interface{}{
		"user_id":   userID,
		"create_at": time.Now(),
	}
	if _, err := db.NamedExec(
		`
			INSERT INTO items (user_id, create_at) 
			VALUES (:user_id, :create_at)
			`,
		info,
	); err != nil {
		log.Printf("ERRER : %v", err)
	}
}

func saveURL(db *sqlx.DB, userID uuid.UUID, k, v string) error {
	info := map[string]interface{}{
		"user_id":      userID,
		"short_url":    k,
		"original_url": v,
	}
	_, err := db.NamedExec(
		`UPDATE items SET short_url=:short_url, original_url=:original_url WHERE user_id=:user_id`,
		info,
	)
	if err != nil {
		log.Printf("ERRER : %v", err)
	}
	return err
}

func checkUser(db *sqlx.DB, id uuid.UUID) (exists bool) {
	err := db.QueryRow(fmt.Sprintf("SELECT exists (%s)", id)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists %v", err)
	}
	return exists
}

func convertToArrayMap(mapURL map[string]string) []map[string]string {
	var urlArray []map[string]string

	for k, v := range mapURL {
		temp := map[string]string{k: v}
		urlArray = append(urlArray, temp)
	}
	return urlArray
}
