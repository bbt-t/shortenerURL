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

	flag.StringVar(&inpFlagParam, "db", "none", `
	"select database: no flag - use map, 'sqlite' - use SQLite, 'pg' - Postgresql, 'redis' - Redis"
	`)
	flag.Parse()

	switch inpFlagParam {
	case "sqlite":
		log.Println("USED SQL")
		db = storage.NewDBSqlite()
	case "pg":
		log.Println("USED PG")
		db = storage.NewDBPostgres()
	case "redis":
		log.Println("USED REDIS")
		db = storage.NewRedisConnect()
	default:
		log.Println("USED MAP")
		db = storage.NewMapDBPlug()
	}

	cfg := configs.NewConfServ()
	h := app.NewHandlerServer(db)

	log.Println("---> RUN SERVER <---")

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.Port), h.Chi))
}
