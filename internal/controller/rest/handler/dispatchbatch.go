package handler

import (
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/pkg"
	"io"
	"log"
	"net/http"
)

func (s ShortenerHandler) buildURLBatch(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for redirecting to original URL.
		get ID from the route  -> search for the original url in DB:
			if not -> 404
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
			item["original_url"] = fmt.Sprintf("%v", pkg.HashShortening([]byte(v)))
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
