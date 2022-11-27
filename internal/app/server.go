package app

import (
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"
	"github.com/bbt-t/shortenerURL/pkg"
	"io"
	"log"
	"net/http"
	"strings"
)

func RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p := strings.Split(r.URL.Path, "/")[1]
		//rUrl, err := st.PullOutUrlRedis(st.RedisClientConnect(), p)
		rUrl, err := st.PullOutUrlSQL(p)

		if err != nil {
			log.Printf("ERROR : %s", err)
			http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)
		}
		w.Header().Set("Location", rUrl)
		w.WriteHeader(307)
		http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)

	case http.MethodPost:
		defer r.Body.Close()

		var value CreateShortURLRequest

		payload, errReadBody := io.ReadAll(r.Body)
		if errReadBody != nil {
			log.Printf("ERROR : %s", errReadBody)
		}

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("ERROR: %s", err)
		}
		fmt.Println(value.URL)

		toHashVar := fmt.Sprintf("%d", pkg.HashShortening([]byte(value.URL)))
		//st.SaveNewUrlRedis(st.RedisClientConnect(), toHashVar, value.URL)
		st.SaveNewUrlSQL(toHashVar, value.URL)

		resp := Resp{
			Result: configs.NewConfig().BaseURL + "/" + toHashVar,
		}
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(res); err != nil {
			log.Printf("ERROR : %s", err)
		}
	}
}
