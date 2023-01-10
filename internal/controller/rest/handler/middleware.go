package handler

import (
	"compress/gzip"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bbt-t/shortenerURL/internal/controller/rest"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func (s ShortenerHandler) GetterSetterAuthJWTCookie(next http.Handler) http.Handler {
	/*
		Cookies-middleware.
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

func (s ShortenerHandler) customGzipCompress(next http.Handler) http.Handler {
	/*
		Custom compress-middleware.
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Del("Content-Length")
		reader, err := gzip.NewReader(r.Body)
		defer reader.Close()

		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		r.Body = gzipReader{
			reader,
			r.Body,
		}
		next.ServeHTTP(w, r)
	})
}

type gzipReader struct {
	*gzip.Reader
	io.Closer
}

func (r gzipReader) Close() error {
	if err := r.Closer.Close(); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}
