package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"time"

	"github.com/bbt-t/shortenerURL/internal/controller/rest"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func (s ShortenerHandler) GetterSetterAuthJWTCookie(next http.Handler) http.Handler {
	/*
		Cookies middleware.
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string

		token, claims, _ := jwtauth.FromContext(r.Context())

		if token != nil && jwt.Validate(token) == nil {
			userID = fmt.Sprintf("%v", claims["user_id"])
		}
		if token == nil || jwt.Validate(token) != nil {

			userUUID, _ := uuid.NewV7()
			s.s.NewUser(userUUID)
			userID = userUUID.String()

			setToken, _ := rest.MakeToken(userID)
			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
				SameSite: http.SameSiteLaxMode,
				//Domain:   pkg.HostOnly(h.cfg.ServerAddress),
				Path: "/",
				// Uncomment below for HTTPS:
				// Secure: true,
				Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
				Value: setToken,
			})
		}
		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
