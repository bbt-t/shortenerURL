package main

import (
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app"
	"github.com/bbt-t/shortenerURL/internal/app/storage"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred : %s", err)
		}
	}()

	db := storage.NewDBSqlite()
	db.CreateSchema()

	cfg := configs.NewConfServ()
	h := app.NewHandlerServer(db)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:8080", cfg.ServerAddress), h.Chi))
}
