package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/go-chi/chi/v5"
)

func (h *ServerHandler) redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for redirecting to original URL.
		get ID from the route  -> search for the original url in DB:
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

	if err := h.store.SaveURL(hashedVal, originalURL); err != nil {
		log.Printf("ERROR : %s", err)
	}

	shortURL := []byte(fmt.Sprintf("%v/%v", h.cfg.BaseURL, hashedVal))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(shortURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
}

func (h *ServerHandler) takeAndSendURLJson(w http.ResponseWriter, r *http.Request) {
	/*
		Comes -> json object {"url": "original_url"}
		Coming out <- response {"result": "shorten_url"}
	*/
	var (
		req  reqURL
		resp respURL
	)

	defer r.Body.Close()
	payload, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(
			w,
			fmt.Sprintf("Incorrent request body: %s", payload),
			http.StatusBadRequest,
		)
		return
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Print(err)
		http.Error(
			w,
			fmt.Sprintf("Impossible unmarshal request : %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	originalURL := strings.TrimSpace(req.URL) // remove spaces
	if originalURL == "" {
		originalURL = string(payload)
	}
	hashedVal := fmt.Sprintf("%d", pkg.HashShortening([]byte(originalURL)))

	if err := h.store.SaveURL(hashedVal, originalURL); err != nil {
		log.Printf("ERROR : %s", err)
	}
	shortURL := fmt.Sprintf("%v/%v", h.cfg.BaseURL, hashedVal)

	resp.URL = shortURL
	result, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}

func (h *ServerHandler) takeAllUrls(w http.ResponseWriter, r *http.Request) {
	var result []map[string]string
	//result, err := h.store.GetAllURL()
	//if err != nil {
	//		http.Error(
	//			w,
	//			fmt.Sprintf("%v", err),
	//			http.StatusNoContent
	//		)
	//		return
	//	}
	jResult, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(jResult); err != nil {
		log.Printf("ERROR : %s", err)
	}
}

func (h *ServerHandler) pingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Ping(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("%v", err),
			http.StatusInternalServerError,
		)
		return
	}
}
