package handler

import (
	"fmt"
	"net/http"
)

func (s ShortenerHandler) pingDB(w http.ResponseWriter, r *http.Request) {
	/*
		Database connection check.
	*/
	if err := s.s.PingDB(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("%v", err),
			http.StatusInternalServerError,
		)
		return
	}
}
