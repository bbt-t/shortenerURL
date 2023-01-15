package storage

import (
	"database/sql"
	"fmt"
	"github.com/bbt-t/shortenerURL/internal/entity"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func createTable(db *sqlx.DB, schema string) {
	db.MustExec(schema /* SQL commands */)
	log.Println("SCHEMA CREATED")
}

func getOriginalURL(db *sqlx.DB, shortURL string) (string, error) {
	var result string

	err := db.Get(
		&result,
		"SELECT original_url FROM items WHERE short_url=$1",
		shortURL,
	)
	return result, err
}

func addNewUser(db *sqlx.DB, userID uuid.UUID) {
	info := map[string]interface{}{
		"user_id":   userID,
		"create_at": time.Now(),
	}
	if _, err := db.NamedExec(
		"INSERT INTO users (user_id, create_at) VALUES (:user_id, :create_at)",
		info,
	); err != nil {
		log.Printf("ERRER : %v", err)
	}
}

func saveURL(db *sqlx.DB, userID uuid.UUID, shortURL, originalURL string) error {
	var check bool
	info := map[string]interface{}{
		"user_id":      userID,
		"original_url": originalURL,
		"short_url":    shortURL,
	}

	if err := db.Get(
		&check,
		"SELECT EXISTS(SELECT 1 FROM items WHERE user_id=$1 AND original_url=$2)",
		userID,
		originalURL,
	); err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists %v", err)
	}
	if check {
		return errHTTPConflict
	}

	_, err := db.NamedExec(
		`
	INSERT INTO items (user_id, original_url, short_url) 
	VALUES (:user_id, :original_url, :short_url)
	`,
		info,
	)
	if err != nil {
		log.Printf("ERRER : %v", err)
	}
	return err
}

func checkUser(db *sqlx.DB, uid uuid.UUID) (exists bool) {
	/*
		Checking if the user exists in the DB.
		param id: user_id
	*/
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)", uid)
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

func saveURLBatch(db *sqlx.DB, uid uuid.UUID, urlBatch []entity.URLBatchInp) error {

	for i, item := range urlBatch {
		temp := strings.Split(item.ShortURL, "/")
		urlBatch[i].ShortURL = temp[len(temp)-1]

		urlBatch[i].UserID = uid

		//uuidURL, _ := uuid.FromString(fmt.Sprintf("%v", urlBatch[i].CorrelationID))
		//urlBatch[i].ID = uuidURL
	}

	query := "INSERT INTO items (user_id, original_url, short_url) VALUES (:user_id, :original_url, :short_url)"

	for _, chunk := range urlBatch {
		if _, err := db.NamedExec(query, &chunk); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
