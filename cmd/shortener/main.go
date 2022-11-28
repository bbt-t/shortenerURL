package main

import (
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

func main() {
	st.AddSchema()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred : %s", err)
		}
	}()
	cfg := configs.NewConfig()

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:8080", cfg.ServerAddress),
		MaxHeaderBytes:    0,
		ReadHeaderTimeout: 1 * time.Second,
	}
	http.HandleFunc("/", app.RedirectToOriginalURL)
	log.Fatal(server.ListenAndServe())
}
