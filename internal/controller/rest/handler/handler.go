package handler

import (
	"github.com/bbt-t/shortenerURL/internal/config"
	"github.com/bbt-t/shortenerURL/internal/controller/rest"
	"github.com/bbt-t/shortenerURL/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	cfg *config.ServerCfg
}

func NewShortenerRoutes(s *usecase.ShortenerService, cfg *config.ServerCfg) *ShortenerHandler {
	return &ShortenerHandler{
		s:   s,
		cfg: cfg,
	}
}

func (s ShortenerHandler) InitRoutes() *chi.Mux {
	/*
		Initialize the server, setting preferences and add routes.
	*/
	//del := NewDeleteHandler(s.s, 10)
	route := chi.NewRouter()
	route.Use(
		middleware.RealIP, // <- (!) Only if a reverse proxy is used (e.g. nginx) (!)
		middleware.Logger,
		middleware.Recoverer,
		// Compress:
		s.customGzipCompress,
		// Working with paths:
		middleware.CleanPath,
		// JWT
		jwtauth.Verifier(rest.TokenAuth),
	)
	// Protected routes:
	route.Group(func(r chi.Router) {
		// Cookie Middleware:
		r.Use(s.GetterSetterAuthJWTCookie)
		// Routes:
		r.Get("/ping", s.pingDB)
		r.Get("/{id}", s.recoverOriginalURL)
		r.Get("/api/user/urls", s.recoverAllOriginalURLByUser)

		r.Delete("/api/user/urls", s.deleteURL) //del.Handle) //s.deleteURL)

		r.Post("/api/shorten/batch", s.buildURLBatch)
		r.Post("/api/shorten", s.composeNewShortURLJson)
		r.Post("/", s.composeNewShortURL)
	})

	return route
}
