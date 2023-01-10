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

func (s ShortenerHandler) Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}
		request.Header.Del("Content-Length")
		reader, err := gzip.NewReader(request.Body)
		if err != nil {
			io.WriteString(writer, err.Error())
			return
		}

		defer reader.Close()

		request.Body = gzipReader{
			reader,
			request.Body,
		}
		log.Println("GZIP MIDDLEWARE")
		next.ServeHTTP(writer, request)
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
