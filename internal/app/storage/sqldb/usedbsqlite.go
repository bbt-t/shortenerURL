package sqldb

import (
	"log"

	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type dbSqlite struct {
	db *sqlx.DB
}

func NewDBSqlite() storage.DBRepo {
	/*
		Initializing the SQLite DB.
		return: DB object
	*/
	db := dbConnect("sqlite3", configs.NewConfSQLite().DBName)
	if err := db.Ping(); err != nil {
		log.Println(err)
	}
	createTable(db, _tableItems /* SQL command */)
	return &dbSqlite{db}
}

func (d *dbSqlite) SaveURL(k, v string) error {
	/*
		Adding info to the DB.
		return: Error or nil
	*/
	err := saveInDB(d.db, k /* id (hashed url) */, v /* original url */)
	return err
}

func (d *dbSqlite) GetURL(k string) (string, error) {
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

func (d *dbSqlite) Ping() error {
	err := d.db.Ping()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("SQLite is READY!")
	}
	return err
}
