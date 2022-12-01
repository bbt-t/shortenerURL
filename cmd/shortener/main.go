package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app"
	"github.com/bbt-t/shortenerURL/internal/app/storage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	/*
		"Application entry point".
		Create table in DB and start the http-server.
	*/
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred : %s", err)
		}
	}()
	var db storage.DBRepo
	var inpFlagParam string

	cfg := configs.NewConfServ()

	flag.StringVar(&inpFlagParam, "db", "none", `
	"select database: no flag - use map, 'sqlite' - use SQLite, 'pg' - Postgresql."
	`)
	flag.Parse()

	switch inpFlagParam {
	case "sqlite":
		db = storage.NewDBSqlite()
	case "pg":
		db = storage.NewDBPostgres()
	default:
		db = storage.NewMapDBPlug()
	}

	h := app.NewHandlerServer(db)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:8080", cfg.ServerAddress), h.Chi))
}
