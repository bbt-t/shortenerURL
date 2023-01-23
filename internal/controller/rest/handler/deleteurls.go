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

	go s.s.DelURLArray(userID, payload)

	w.WriteHeader(http.StatusAccepted)
}
