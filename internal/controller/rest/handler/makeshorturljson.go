package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/gofrs/uuid"
)

func (s ShortenerHandler) composeNewShortURLJson(w http.ResponseWriter, r *http.Request) {
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

	shortURL := fmt.Sprintf("%d", pkg.HashShortening([]byte(originalURL)))

	temp := r.Context().Value("user_id")
	userID, _ := uuid.FromString(fmt.Sprintf("%v", temp))

	errSaveURL := s.s.SaveShortURL(userID, shortURL, originalURL)
	resp.URL = fmt.Sprintf("%v/%v", s.cfg.BaseURL, shortURL)

	result, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")

	if errSaveURL != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
