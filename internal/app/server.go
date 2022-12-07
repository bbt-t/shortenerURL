package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
}

func NewHandlerServer(s st.DBRepo) *ServerHandler {
	/*
		Initialize the server and add routes.
	*/
	router := chi.NewRouter()
	h := ServerHandler{
		Chi:   router,
		store: s,
	}

	h.Chi.Get("/{id}", h.redirectToOriginalURL)
	h.Chi.Post("/", h.takeAndSendURL)
	return &h
}

func Start(inpFlagParam string) {
	/*
		Get param, choice of storage to use
		(if the selected storage is not available, then the MAP is selected)
		and start the http-server.
	*/
	var db st.DBRepo

	if inpFlagParam != "redis" {
		db = st.NewSQLDatabase(inpFlagParam /* flag for choice DB */)
	} else {
		db = st.NewRedisConnect()
	}
	if nil == db {
		db = st.NewMapDBPlug()
		log.Println("--->>> SWITCH TO MAP")
	}

	cfg := configs.NewConfServ()
	h := NewHandlerServer(db)
	log.Println("---> RUN SERVER <---")
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.Port),
			h.Chi,
		),
	)
}
