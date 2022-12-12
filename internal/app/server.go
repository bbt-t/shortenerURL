package app

import (
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
	cfg   configs.ServerCfg
}

func NewHandlerServer(s st.DBRepo, cfg configs.ServerCfg) *ServerHandler {
	/*
		Initialize the server and add routes.
	*/
	allowedCharsets := []string{"UTF-8", "Latin-1", ""}

	router := chi.NewRouter()
	h := ServerHandler{
		Chi:   router,
		store: s,
		cfg:   cfg,
	}

	h.Chi.Use(middleware.Logger)
	h.Chi.Use(middleware.Recoverer)

	h.Chi.Use(middleware.CleanPath)
	h.Chi.Use(middleware.RedirectSlashes)

	h.Chi.Use(middleware.ContentCharset(allowedCharsets...))
	h.Chi.Use(middleware.AllowContentType("application/json", "text/plain"))
	h.Chi.Use(middleware.AllowContentEncoding("deflate", "gzip"))

	h.Chi.Use(middleware.Compress(5, "application/json", "text/plain"))

	h.Chi.Get("/{id}", h.redirectToOriginalURL)
	h.Chi.Post("/api/shorten", h.takeAndSendURLJson)
	h.Chi.Post("/", h.takeAndSendURL)

	return &h
}

func Start(cfg *configs.ServerCfg) {
	/*
		Get param, choice of storage to use
		(if the selected storage is not available, then the MAP is selected)
		and start the http-server.
	*/
	var db st.DBRepo

	if cfg.FilePath != "" {
		log.Println("WITH FILE STORAGE --->>>")
		db = st.NewFileDB(cfg.FilePath)
	} else {
		if cfg.UseDB != "redis" {
			db = st.NewSQLDatabase(cfg.UseDB)
		} else {
			db = st.NewRedisConnect()
		}
		if nil == db {
			db = st.NewMapDBPlug()
			log.Println("--->>> SWITCH TO MAP")
		}
	}

	h := NewHandlerServer(db, *cfg)
	log.Println("---> RUN SERVER <---")
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, h.Chi))
}
