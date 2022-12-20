package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/go-chi/jwtauth/v5"
)

func (h *ServerHandler) singJWTCookie(w http.ResponseWriter, r *http.Request) {
	var userName, userPassword string

	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			log.Print(err)
			http.Error(
				w,
				fmt.Sprintf("Error : %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		userName = r.PostForm.Get("username")
		userPassword = r.PostForm.Get("password")
	case "application/json":
		var userInfo loginIn

		defer r.Body.Close()
		payload, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(payload, &userInfo); err != nil {
			log.Print(err)
			http.Error(
				w,
				fmt.Sprintf("Impossible unmarshal request : %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		userName = userInfo.UserName
		userPassword = userInfo.Password
	}
	if userName == "" || userPassword == "" {
		http.Error(
			w,
			"Missing username or password.",
			http.StatusBadRequest,
		)
		return
	}

	token, _ := makeToken(userName)
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Domain:   pkg.HostOnly(h.cfg.ServerAddress),
		Path:     "/admin",
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ServerHandler) adminAuth(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	if _, err := w.Write([]byte(fmt.Sprintf("Hi %v", claims["username"]))); err != nil {
		log.Printf("ERROR : %s", err)
	}
}

// TODO:
// 		1. сделать общий интерфейс для json структур с методами маршал, анмаршал.
//  	2. сделать регистрацию в бд
//		3. проверку логин/пасс в бд через encrypt
