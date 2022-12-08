package app

import (
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/pkg"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *ServerHandler) takeAndSendURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for getting URL to shortened.
		Received, run through the HASH-func and write (hash, original url)
		to the DB and (hash only) response Body, sent response.
	*/
	cfg := configs.NewConfServ()

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
	//if !pkg.URLValidation(originalURL) {
	//	http.Error(
	//		w,
	//		fmt.Sprintf("Incorrent URL: %s", payload),
	//		http.StatusBadRequest,
	//	)
	//}
	hashedVal := fmt.Sprintf("%d", pkg.HashShortening([]byte(originalURL)))

	if err := h.store.SaveURL(hashedVal, originalURL); err != nil {
		log.Printf("ERROR : %s", err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	shortURL := []byte(
		fmt.Sprintf(
			"http://%v:%v/%v", cfg.ServerAddress, cfg.Port, hashedVal),
	)
	if _, err := w.Write(shortURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
