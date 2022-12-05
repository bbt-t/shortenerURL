package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	"github.com/bbt-t/shortenerURL/internal/app/storage/nosqldb"
	"github.com/bbt-t/shortenerURL/internal/app/storage/sqldb"
	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
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

	h.Chi.Get("/{id}", h.RedirectToOriginalURL)
	h.Chi.Post("/", h.TakeAndSendURL)
	return &h
}

func (h *ServerHandler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for redirecting to original URL.
		Get ID from the route  -> search for the original url in DB:
			if it's found -> redirect
			if not -> 404
	*/
	if originalURL, err := h.store.GetURL(chi.URLParam(r, "id")); err != nil {
		log.Printf("ERROR : %s", err)
		http.NotFound(w, r)
	} else {
		w.Header().Set("Location", originalURL)
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}

func (h *ServerHandler) TakeAndSendURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for getting URL to shortened.
		Received, run through the HASH-func and write (hash, original url)
		to the DB and (hash only) response Body, sent response.
	*/
	var value CreateShortURLRequest
	var shortURL []byte
	cfg := configs.NewConfServ()

	defer r.Body.Close()
	payload, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		log.Printf("ERROR : %s", errReadBody)
	}
	if err := json.Unmarshal(payload, &value); err != nil {
		log.Printf("ERROR: %s", err)
	}

	if pkg.URLValidation(value.URL) == true {
		toHashVar := fmt.Sprintf("%d", pkg.HashShortening([]byte(value.URL)))

		if err := h.store.SaveURL(toHashVar, value.URL); err != nil {
			log.Printf("ERROR : %s", err)
		}

		shortURL = []byte(
			fmt.Sprintf(
				"http://%s:%s/%s", cfg.ServerAddress, cfg.Port, toHashVar),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(shortURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
}

func Start(inpFlagParam string) {
	/*
		Get param, choice of storage to use
		(if the selected storage is not available, then the MAP is selected)
		and start the http-server.
	*/
	var db st.DBRepo

	switch inpFlagParam {
	case "sqlite":
		log.Println("USED SQL")
		db = sqldb.NewDBSqlite()
	case "pg":
		log.Println("USED PG")
		db = sqldb.NewDBPostgres()
	case "redis":
		log.Println("USED REDIS")
		db = nosqldb.NewRedisConnect()
	default:
		log.Println("USED MAP")
		db = st.NewMapDBPlug()
	}

	if err := db.Ping(); err != nil {
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
