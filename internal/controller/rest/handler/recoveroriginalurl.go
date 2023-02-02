package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s ShortenerHandler) recoverOriginalURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for redirecting to original URL.
		get ID from the route  -> search for the original url in DB:
			if not -> 404
			if status deleted -> 410
	*/
	if originalURL, err := s.s.GetOriginalURL(chi.URLParam(r, "id")); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.NotFound(w, r)
	} else {
		w.Header().Set("Location", originalURL)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
