package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"net/http"

	"github.com/bbt-t/shortenerURL/pkg"
)

func (s ShortenerHandler) buildURLBatch(w http.ResponseWriter, r *http.Request) {
	/*
		Accepts multiple URLs in the request body to shorten,
		changes "original_url" to "short_url".
	*/
	var urlBatch []map[string]string

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
	if err := json.Unmarshal(payload, &urlBatch); err != nil {
		log.Print(err)
		http.Error(
			w,
			fmt.Sprintf("Impossible unmarshal request : %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	temp := r.Context().Value("user_id")
	userID, _ := uuid.FromString(fmt.Sprintf("%v", temp))

	for _, item := range urlBatch {
		for _, v := range item {
			shortURL := fmt.Sprintf("%v", pkg.HashShortening([]byte(v)))

			if err := s.s.SaveShortURL(userID, shortURL, item["original_url"]); err != nil {
				log.Print(err)
				http.Error(
					w,
					fmt.Sprintf("Impossible unmarshal request : %s", err),
					http.StatusInternalServerError,
				)
				return
			}

			item["short_url"] = shortURL
			delete(item, "original_url")
		}
	}

	result, err := json.Marshal(urlBatch)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
