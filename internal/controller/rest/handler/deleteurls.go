package handler

import (
	"fmt"
	"io"
	"net/http"
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

	//go h.s.DelURLArray()
	w.WriteHeader(http.StatusAccepted)
}
