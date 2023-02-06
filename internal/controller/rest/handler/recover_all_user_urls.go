package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s ShortenerHandler) recoverAllOriginalURLByUser(w http.ResponseWriter, r *http.Request) {
	/*
		Get user id -> go to db -> write to json -> response
	*/
	userID, _ := uuid.FromString(fmt.Sprintf("%v", r.Context().Value("user_id")))

	allURL, errGetURL := s.s.GetURLArrayByUser(userID, s.cfg.BaseURL)
	if errGetURL != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, errJSON := json.Marshal(allURL)
	if errJSON != nil {
		log.Println(errJSON)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
