package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

func createTable(ctx context.Context, db *pgxpool.Pool, schema string) {
	/*
		Executes SQL commands at startup.
		param schema: commands
	*/
	if _, err := db.Exec(ctx, schema /* SQL commands */); err != nil {
		log.Fatal(err)
	}
	log.Println("SCHEMA CREATED")
}

func getOriginalURL(db *pgxpool.Pool, shortURL string) (string, error) {
	/*
		Makes a request (selection of the original url by short) to the DB.
		return: original url or error
	*/
	var result entity.CheckURL
	ctx := context.Background()
	err := db.QueryRow(
		ctx,
		"SELECT original_url, deleted FROM items WHERE short_url=$1",
		shortURL,
	).Scan(&result.OriginalURL, &result.Deleted)
	if result.Deleted {
		return "", errDeleted
	}
	return result.OriginalURL, err
}

func addNewUser(db *pgxpool.Pool, userID uuid.UUID) {
	/*
		Adds a new user to the DB.
		param userID: UUID issued when receiving the cookie (middleware)
	*/
	ctx := context.Background()
	if _, err := db.Exec(ctx, "INSERT INTO users (user_id) VALUES ($1)", userID); err != nil {
		log.Printf("ERRER : %+v", err)
	}
}

func saveURL(db *pgxpool.Pool, userID uuid.UUID, shortURL, originalURL string) error {
	/*
		Adds short url to DB.
	*/
	ctx := context.Background()

	var check bool
	if err := db.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM items WHERE user_id=$1 AND original_url=$2)",
		userID,
		originalURL,
	).Scan(&check); err != nil && err != pgx.ErrNoRows { // ??? errors.Is()
		log.Printf("error checking if row exists %+v", err)
	}
	if check {
		return errHTTPConflict
	}

	_, err := db.Exec(
		ctx,
		"INSERT INTO items (user_id, original_url, short_url) values($1, $2, $3)",
		userID,
		originalURL,
		shortURL,
	)

	return err
}

func getOriginalURLArray(db *pgxpool.Pool, userID uuid.UUID, baseURL string) ([]map[string]string, error) {
	var resultStructs []*entity.URLs // *???
	ctx := context.Background()
	rows, err := db.Query(
		ctx,
		"SELECT original_url, short_url FROM items WHERE user_id=$1 and deleted is false",
		userID,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		url := &entity.URLs{} // &???
		if err = rows.Scan(&url.OriginalURL, &url.ShortURL); err == nil {
			resultStructs = append(resultStructs, url)
		}
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

func checkUser(db *pgxpool.Pool, uid uuid.UUID) (exists bool) {
	/*
		Checking if the user exists in the DB.
		param uid: UUID issued when receiving the cookie (middleware)
	*/
	ctx := context.Background()
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)", uid).Scan(&exists)
	if !errors.Is(err, pgx.ErrNoRows) {
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

func saveURLBatch(ctx context.Context, db *pgxpool.Pool, uid uuid.UUID, urlBatch []entity.URLBatchInp) error {
	for i, item := range urlBatch {
		temp := strings.Split(item.ShortURL, "/")
		urlBatch[i].ShortURL = temp[len(temp)-1]
		urlBatch[i].UserID = uid
	}

	batch := &pgx.Batch{}

	for _, url := range urlBatch {
		batch.Queue(
			"INSERT INTO items (user_id, original_url, short_url) VALUES($1, $2, $3)",
			url.UserID,
			url.OriginalURL,
			url.ShortURL,
		)
	}

	query := db.SendBatch(ctx, batch)
	defer query.Close()

	_, err := query.Exec()

	return err
}

func deleteURLArray(ctx context.Context, db *pgxpool.Pool, uid uuid.UUID, inpURLs []string) error {
	_, err := db.Exec(
		ctx,
		"UPDATE items SET deleted=true WHERE user_id = $1 AND short_url = any($2::text[])",
		uid,
		pq.Array(inpURLs),
	)
	if err != nil {
		log.Println(err)
	}
	return err
}
