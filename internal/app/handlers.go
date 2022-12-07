package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/go-chi/chi/v5"
)

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
	var shortURLRequest ResShortURL
	var shortURL []byte
	cfg := configs.NewConfServ()

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	if err := json.Unmarshal(payload, &shortURLRequest); err != nil {
		log.Printf("ERROR: %s", err)
	}

	if pkg.URLValidation(shortURLRequest.URL) {
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
