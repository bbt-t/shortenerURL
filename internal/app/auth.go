package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

const _secret = "<jwt-secret>"

var _tokenAuth *jwtauth.JWTAuth

func init() {
	_tokenAuth = jwtauth.New("HS256", []byte(_secret), nil)
}

func makeToken(name string) (string, error) {
	_, tokenString, err := _tokenAuth.Encode(map[string]interface{}{"username": name})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *ServerHandler) singJWT(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(
			w,
			fmt.Sprintf("Error : %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	userName := r.PostForm.Get("username")
	userPassword := r.PostForm.Get("password")

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
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

	//_, claims, _ := jwtauth.FromContext(r.Context())
	//if _, err := w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"]))); err != nil {
	//	log.Printf("ERROR : %s", err)
	//}
}
