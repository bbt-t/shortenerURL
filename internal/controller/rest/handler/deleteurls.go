package handler

import (
	"context"
	"fmt"
	"github.com/bbt-t/shortenerURL/pkg"
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

	ctx := context.Background()

	ids := pkg.ConvertStrToSlice(string(payload))
	go s.s.DelURLArray(ctx, userID, ids)

	w.WriteHeader(http.StatusAccepted)
}
