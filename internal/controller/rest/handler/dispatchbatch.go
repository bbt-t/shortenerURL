package handler

import (
	"encoding/json"
	"fmt"
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

	for _, item := range urlBatch {
		for _, v := range item {
			delete(item, "original_url")
			item["short_url"] = fmt.Sprintf("%s/%v", s.cfg.BaseURL, pkg.HashShortening([]byte(v)))
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
