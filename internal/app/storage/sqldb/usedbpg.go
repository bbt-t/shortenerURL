package sqldb

import (
	"log"

	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type dbPostgres struct {
	db *sqlx.DB
}

func NewDBPostgres() storage.DBRepo {
	/*
		Initializing the Postgresql DB.
		return: DB object
	*/
	db := dbConnect("postgres", configs.NewConfPG().DBUrl)

	err := db.Ping()
	if err != nil {
		log.Println(err)
	}
	return &dbPostgres{db}
}

func (d dbPostgres) SaveURL(k, v string) error {
	return nil
}

func (d dbPostgres) GetURL(k string) (string, error) {
	result, err := getInfo(d.db, k)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (d dbPostgres) Ping() error {
	err := d.db.Ping()
	if err != nil {
		log.Println(err)
	}
	return err
}
