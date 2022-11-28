package app

import (
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	"github.com/bbt-t/shortenerURL/pkg"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strings"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
}

func NewHandlerServer(s st.DBRepo) *ServerHandler {
	router := chi.NewRouter()
	h := ServerHandler{
		Chi:   router,
		store: s,
	}

	h.Chi.Get("/{id}", h.RedirectToOriginalURL)
	h.Chi.Post("/", h.TakeAndSendUrl)
	return &h
}

func (h ServerHandler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1]
	rUrl, err := h.store.GetURL(p)
	if err != nil {
		log.Printf("ERROR : %s", err)
		http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)
	}
	fmt.Println(rUrl)
	w.Header().Set("Location", rUrl)
	http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)
}

func (h ServerHandler) TakeAndSendUrl(w http.ResponseWriter, r *http.Request) {
	var value CreateShortURLRequest

	defer r.Body.Close()
	payload, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		log.Printf("ERROR : %s", errReadBody)
	}
	if err := json.Unmarshal(payload, &value); err != nil {
		log.Printf("ERROR: %s", err)
	}
	toHashVar := fmt.Sprintf("%d", pkg.HashShortening([]byte(value.URL)))

	if err := h.store.SaveURL(toHashVar, value.URL); err != nil {
		log.Printf("ERROR : %s", err)
	}
	//resp := Resp{
	//	Result: configs.NewConfServ().BaseURL + "/" + toHashVar,
	//}
	//res, err := json.Marshal(resp)
	//if err != nil {
	//	panic(err)
	//}
	resp := []byte(configs.NewConfServ().BaseURL + "/" + toHashVar)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(resp); err != nil {
		log.Printf("ERROR : %s", err)
	}
}
