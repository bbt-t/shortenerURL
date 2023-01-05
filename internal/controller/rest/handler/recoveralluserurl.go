package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

func (s ShortenerHandler) recoverAllOriginalURLByUser(w http.ResponseWriter, r *http.Request) {
	/*
		получить ИД пользователя -> сходить в БД -> записать в json -> выдать
	*/
	temp := r.Context().Value("user_id")
	userID, _ := uuid.FromString(fmt.Sprintf("%v", temp))

	allURL, errGetURL := s.s.GetURLArrayByUser(userID)
	result, errJson := json.Marshal(allURL)
	if errJson != nil {
		log.Println(errJson)
	}

	w.Header().Set("Content-Type", "application/json")
	if errGetURL != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if _, err := w.Write(result); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
