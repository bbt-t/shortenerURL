package sqldb

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func dbConnect(dbName, info string) *sqlx.DB {
	/*
		Connected to DataBase.
		param dbName: driver name (sqlite or postgres)
		param info: info for conn. (sqlite - name db, postgres - db url)
	*/
	db, err := sqlx.Connect(dbName, info)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	return db
}

func createTable(db *sqlx.DB, schema string) {
	db.MustExec(schema)
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
