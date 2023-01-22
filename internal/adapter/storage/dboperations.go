package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bbt-t/shortenerURL/internal/entity"
	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func createTable(db *sqlx.DB, schema string) {
	/*
		Executes SQL commands at startup.
		param schema: commands
	*/
	db.MustExec(schema /* SQL commands */)
	log.Println("SCHEMA CREATED")
}

func getOriginalURL(db *sqlx.DB, shortURL string) (string, error) {
	/*
		Makes a request (selection of the original url by short) to the DB.
		return: original url or error
	*/
	var result entity.CheckURL

	err := db.Get(
		&result,
		"SELECT original_url, deleted FROM items WHERE short_url=$1",
		shortURL,
	)
	if result.Deleted {
		return "", errDeleted
	}
	return result.OriginalURL, err
}

func addNewUser(db *sqlx.DB, userID uuid.UUID) {
	/*
		Adds a new user to the DB.
		param userID: UUID issued when receiving the cookie (middleware)
	*/
	info := map[string]interface{}{
		"user_id":   userID,
		"create_at": time.Now(),
	}
	if _, err := db.NamedExec(
		"INSERT INTO users (user_id, create_at) VALUES (:user_id, :create_at)",
		info,
	); err != nil {
		log.Printf("ERRER : %+v", err)
	}
}

func saveURL(db *sqlx.DB, userID uuid.UUID, shortURL, originalURL string) error {
	/*
		Adds short url to DB.
	*/
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
		log.Printf("error checking if row exists %+v", err)
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
		log.Printf("ERRER : %+v", err)
	}
	return err
}

func getOriginalURLArray(db *sqlx.DB, userID uuid.UUID, baseURL string) ([]map[string]string, error) {
	var resultStructs []entity.URLs

	err := db.Select(&resultStructs, "SELECT original_url, short_url FROM items WHERE user_id=$1", userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(resultStructs) == 0 {
		return nil, errDBEmpty
	}
	urlArray := make([]map[string]string, len(resultStructs))

	for _, item := range resultStructs {
		temp := make(map[string]string, 2)

		data, _ := json.Marshal(item)
		_ = json.Unmarshal(data, &temp)

		temp["short_url"] = fmt.Sprintf("%s/%s", baseURL, temp["short_url"])
		urlArray = append(urlArray, temp)
	}

	return urlArray, nil
}

func checkUser(db *sqlx.DB, uid uuid.UUID) (exists bool) {
	/*
		Checking if the user exists in the DB.
		param uid: UUID issued when receiving the cookie (middleware)
	*/
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)", uid)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists %+v", err)
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

func saveURLBatch(ctx context.Context, db *sqlx.DB, uid uuid.UUID, urlBatch []entity.URLBatchInp) error {
	for i, item := range urlBatch {
		temp := strings.Split(item.ShortURL, "/")
		urlBatch[i].ShortURL = temp[len(temp)-1]
		urlBatch[i].UserID = uid
	}

	query := `
			INSERT INTO items (user_id, original_url, short_url) 
			VALUES (:user_id, :original_url, :short_url) ON CONFLICT DO NOTHING
			`
	if rows, err := db.NamedQueryContext(ctx, query, urlBatch); rows.Err() != nil {
		return err
	}

	return nil
}

type ForTestDel struct {
	Deleted  bool      `db:"deleted"`
	ShortURL string    `db:"short_url"`
	UserID   uuid.UUID `db:"user_id"`
}

func deleteURLArray(db *sqlx.DB, uid uuid.UUID, inpJSON []byte) error {
	inpURLs := pkg.ConvertStrToSlice(string(inpJSON))

	//var info []map[string]interface{}

	//for _, v := range inpURLs {
	//	info = append(info, map[string]interface{}{
	//		"deleted":   true,
	//		"short_url": v,
	//		"user_id":   uid,
	//	})
	//}

	query := "UPDATE items SET deleted=:deleted WHERE user_id=:user_id AND short_url=:short_url"

	for _, v := range inpURLs {
		if _, err := db.NamedExec(query, map[string]interface{}{
			"deleted":   true,
			"short_url": v,
			"user_id":   uid,
		}); err != nil {
			return err
		}
	}

	return nil
}
