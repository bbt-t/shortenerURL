package storage

import (
	"fmt"
	"log"

	"github.com/bbt-t/shortenerURL/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqlDatabase struct {
	// Single struct for SQLite and PG.
	db *sqlx.DB
}

func NewSQLDatabase(nameDB, dbURL string) DBRepo {
	/*
		Selects sql-db and initializing. Create tables.
		param nameDB: received parameter (flag) to select db
		return: db-object or nil
	*/
	var param string

	switch nameDB {
	case "sqlite":
		nameDB = fmt.Sprintf("%s3", nameDB)
		param = configs.NewConfSQLite().DBName
	case "pg":
		nameDB = "postgres"
		param = configs.NewConfPG(dbURL).DBUrl
	}

	db, err := sqlx.Connect(nameDB, param)
	if err != nil {
		log.Print(errDBNotSelected /* custom error */)
		return nil
	}

	createTable(db, _tableItems /* SQL command */)
	return &sqlDatabase{db: db}
}

func (d *sqlDatabase) SaveURL(k, v string) error {
	/*
		Adding info to the DB.
		return: Error or nil
	*/
	err := saveInDB(d.db, k /* id (hashed url) */, v /* original url */)
	return err
}

func (d *sqlDatabase) GetURL(k string) (string, error) {
	/*
		Search for info by ID.
		param k: id by which we search in the DB
		return: result (or "") and error (or nil)
	*/
	result, err := getInfo(d.db, k /* id (hashed url) */)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *sqlDatabase) Ping() error {
	err := d.db.Ping()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("SQL DB is READY!")
	}
	return err
}
