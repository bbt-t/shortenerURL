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
			INSERT INTO users (user_id, create_at) 
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
		"original_url": v,
		"short_url":    k,
		"create_at":    time.Now(),
	}
	var check bool
	if err := db.Get(&check, "SELECT EXISTS(SELECT 1 FROM items WHERE user_id=$1 AND original_url=$2)", userID, v); err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists %v", err)
	}
	if check {
		return errHTTPConflict
	}

	_, err := db.NamedExec(`INSERT INTO items (user_id, original_url, short_url, create_at) VALUES (:user_id, :original_url, :short_url, :create_at)`,
		info,
	)
	if err != nil {
		log.Printf("ERRER : %v", err)
	}
	return err
}

func checkUser(db *sqlx.DB, id uuid.UUID) (exists bool) {
	/*
		Checking if the user exists in the DB.
		param id: user_id
	*/
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)", id)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists %v", err)
	}
	return exists
}

func convertToArrayMap(mapURL map[string]string, baseURL string) []map[string]string {
	var urlArray []map[string]string

	for k, v := range mapURL {
		temp := map[string]string{
			"short_url":    fmt.Sprintf("%s/%s", baseURL, k),
			"original_url": v,
		}
		urlArray = append(urlArray, temp)
	}
	return urlArray
}
