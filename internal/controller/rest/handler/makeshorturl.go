package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/gofrs/uuid"
)

func (s ShortenerHandler) composeNewShortURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for getting URL to shortened.
		Received, run through the HASH-func and write (hash, original url)
		to the DB and (hash only) response Body, sent response.
	*/
	defer r.Body.Close()
	payload, _ := io.ReadAll(r.Body)
	query, err := url.ParseQuery(string(payload))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Incorrent request body: %s", payload),
			http.StatusBadRequest,
		)
		return
	}
	originalURL := strings.TrimSpace(query.Get("url")) // remove spaces
	if originalURL == "" {
		originalURL = string(payload)
	}

	hashedVal := fmt.Sprintf("%d", pkg.HashShortening([]byte(originalURL)))

	temp := r.Context().Value("user_id")
	userID, _ := uuid.FromString(fmt.Sprintf("%v", temp))

	if err := s.s.SaveShortURL(userID, hashedVal, originalURL); err != nil {
		log.Printf("ERROR : %s", err)
	}

	shortURL := []byte(fmt.Sprintf("%v/%v", s.cfg.BaseURL, hashedVal))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(shortURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
