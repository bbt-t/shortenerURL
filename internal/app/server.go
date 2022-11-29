package app

import (
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	"github.com/bbt-t/shortenerURL/pkg"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
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

	h.Chi.Get("/{id}", h.RedirectToOriginalURL)
	h.Chi.Post("/", h.TakeAndSendUrl)
	return &h
}

func (h ServerHandler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for redirecting to original URL.
		Get ID from the route  -> search for the original url in DB:
			if it's found -> redirect
			if not -> 404
	*/

	//p := strings.Split(r.URL.Path, "/")[1]
	if originalURL, err := h.store.GetURL(chi.URLParam(r, "id")); err != nil {
		log.Printf("ERROR : %s", err)
		http.NotFound(w, r)
	} else {
		w.Header().Set("Location", originalURL)
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}

func (h ServerHandler) TakeAndSendUrl(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for getting URL to shortened.
		Received, run through the HASH-func and write (hash, original url)
		to the DB and (hash only) response Body, sent response.
	*/
	var value CreateShortURLRequest

	defer r.Body.Close()
	payload, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		log.Printf("ERROR : %s", errReadBody)
	}
	if err := json.Unmarshal(payload, &value); err != nil {
		log.Printf("ERROR: %s", err)
	}
	toHashVar := fmt.Sprintf("%d", pkg.HashShortening([]byte(value.URL)))

	if err := h.store.SaveURL(toHashVar, value.URL); err != nil {
		log.Printf("ERROR : %s", err)
	}

	shortURL := []byte(configs.NewConfServ().BaseURL + "/" + toHashVar)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(shortURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
