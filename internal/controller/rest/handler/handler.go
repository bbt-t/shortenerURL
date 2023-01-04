package handler

import (
	"time"

	"github.com/bbt-t/shortenerURL/internal/config"
	"github.com/bbt-t/shortenerURL/internal/controller/rest"
	"github.com/bbt-t/shortenerURL/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

type reqURL struct {
	URL string `json:"url"`
}

type respURL struct {
	URL string `json:"result"`
}

type ShortenerHandler struct {
	s   *usecase.ShortenerService
	cfg *configs.ServerCfg
}

func NewShortenerRoutes(s *usecase.ShortenerService, cfg *configs.ServerCfg) *ShortenerHandler {
	return &ShortenerHandler{
		s:   s,
		cfg: cfg,
	}
}

func (s ShortenerHandler) InitRoutes() *chi.Mux {
	route := chi.NewRouter()

	/*
		Initialize the server, setting preferences and add routes.
	*/
	allowedCharsets, allowContentTypes :=
		[]string{
			"UTF-8",
			"Latin-1",
			"",
		},
		[]string{
			"application/json",
			"text/plain",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
		}

	route.Use(
		//middleware.RealIP, // <- (!) Only if a reverse proxy is used (e.g. nginx) (!)
		middleware.Logger,
		middleware.Recoverer,
		// Working with paths:
		middleware.CleanPath,
		middleware.RedirectSlashes,
		// Throttle:
		middleware.ThrottleBacklog(10, 50, time.Second*10),
		httprate.LimitByIP(100, 1*time.Minute),
		// Allowed content:
		middleware.ContentCharset(allowedCharsets... /* list unpacking */),
		middleware.AllowContentType(allowContentTypes... /* list unpacking */),
		// Compress:
		middleware.AllowContentEncoding("gzip"),
		middleware.Compress(5, "application/json", "text/plain"),
		// JWT
		jwtauth.Verifier(rest.TokenAuth),
	)

	// Protected routes:
	route.Group(func(r chi.Router) {
		r.Use(s.GetterSetterAuthJWTCookie)

		r.Get("/ping", s.pingDB)
		r.Get("/{id}", s.recoverOriginalURL)
		//r.Post("/api/shorten/batch", ...) // <- for inc12
		r.Post("/api/shorten", s.composeNewShortURLJson)
		r.Post("/", s.composeNewShortURL)
	})

	return route
}
