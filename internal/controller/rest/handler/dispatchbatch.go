package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/internal/entity"
	"github.com/bbt-t/shortenerURL/pkg"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

func (s ShortenerHandler) buildURLBatch(w http.ResponseWriter, r *http.Request) {
	/*
		Accepts multiple URLs in the request body to shorten,
		changes "original_url" to "short_url".
	*/
	var urlBatchForSave []entity.URLBatchInp
	var urlBatchForSend []entity.URLBatch

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
	if err := json.Unmarshal(payload, &urlBatchForSave); err != nil {
		log.Print(err)
		http.Error(
			w,
			fmt.Sprintf("Impossible unmarshal request : %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	for i, item := range urlBatchForSave {
		shortURL := fmt.Sprintf("%v", pkg.HashShortening([]byte(item.OriginalURL)))
		urlBatchForSave[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL)
	}

	userID, _ := uuid.FromString(fmt.Sprintf("%v", r.Context().Value("user_id")))

	// must do deepcopy (!):
	copySt := append(make([]entity.URLBatchInp, 0, len(urlBatchForSave)), urlBatchForSave...)

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	if err := s.s.SaveURLArray(ctx, userID, copySt); err != nil {
		log.Print(err)
		http.Error(
			w,
			"!!! too much lines !!!",
			http.StatusRequestEntityTooLarge,
		)
		return
	}

	temp, _ := json.Marshal(urlBatchForSave)
	_ = json.Unmarshal(temp, &urlBatchForSend)

	result, err := json.Marshal(urlBatchForSend)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
