package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s ShortenerHandler) deleteURL(w http.ResponseWriter, r *http.Request) {
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
	userID, _ := uuid.FromString(fmt.Sprintf("%v", r.Context().Value("user_id")))

	if err := s.s.DelURLArray(userID, payload); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Invalid prameters: %s", payload),
			http.StatusConflict,
		)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
