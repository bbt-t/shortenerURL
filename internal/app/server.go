package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
}

type createShortURLRequest struct {
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

	h.Chi.Get("/{id}", h.redirectToOriginalURL)
	h.Chi.Post("/", h.takeAndSendURL)
	return &h
}

func (h *ServerHandler) redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
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

func (h *ServerHandler) takeAndSendURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for getting URL to shortened.
		Received, run through the HASH-func and write (hash, original url)
		to the DB and (hash only) response Body, sent response.
	*/
	var shortURLRequest createShortURLRequest
	var shortURL []byte
	cfg := configs.NewConfServ()

	defer r.Body.Close()
	payload, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		log.Printf("ERROR : %s", errReadBody)
	}
	if err := json.Unmarshal(payload, &shortURLRequest); err != nil {
		log.Printf("ERROR: %s", err)
	}

	if pkg.URLValidation(shortURLRequest.URL) == true {
		hashedVal := fmt.Sprintf("%d", pkg.HashShortening([]byte(shortURLRequest.URL)))

		if err := h.store.SaveURL(hashedVal, shortURLRequest.URL); err != nil {
			log.Printf("ERROR : %s", err)
		}

		shortURL = []byte(
			fmt.Sprintf(
				"http://%s:%s/%s", cfg.ServerAddress, cfg.Port, hashedVal),
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
