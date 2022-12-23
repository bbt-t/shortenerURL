package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbt-t/shortenerURL/configs"
	st "github.com/bbt-t/shortenerURL/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

type ServerHandler struct {
	Chi   *chi.Mux
	store st.DBRepo
	cfg   configs.ServerCfg
}

func NewHandlerServer(s st.DBRepo, cfg configs.ServerCfg) *ServerHandler {
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

	h := ServerHandler{
		Chi:   chi.NewRouter(),
		store: s,
		cfg:   cfg,
	}

	h.Chi.Use(
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
	)

	// Protected routes:
	h.Chi.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(_tokenAuth))

		r.Post("/login", h.singJWTCookie)
	})

	h.Chi.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(_tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", h.adminAuth)
	})

	// Public routes:
	h.Chi.Group(func(r chi.Router) {
		r.Get("/ping", h.pingDB)
		r.Get("/{id}", h.recoverOriginalURL)
		r.Get("/api/user/urls", h.takeAllUrls)

		//r.Post("/api/shorten/batch", ...) // <- for inc12
		r.Post("/api/shorten", h.composeNewShortURLJson)
		r.Post("/", h.composeNewShortURL)
	})

	return &h
}

func Start(cfg *configs.ServerCfg) {
	/*
		Get param, choice of storage to use
		(if the selected storage is not available, then the MAP is selected)
		and start the http-server.
	*/
	var db st.DBRepo
	// Database selection to use:
	if cfg.FilePath != "" {
		log.Println("WITH FILE STORAGE -->")
		db = st.NewFileDB(cfg.FilePath)
	} else {
		if cfg.UseDB != "redis" {
			db = st.NewSQLDatabase(cfg.UseDB, cfg.DBConnectURL)
		} else {
			db = st.NewRedisConnect()
		}
		if nil == db {
			db = st.NewMapDBPlug()
			log.Println("--> SWITCH TO MAP")
		}
	}

	h := NewHandlerServer(db, *cfg)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: h.Chi,
	}
	// Graceful shutdown:
	// Taken from Chi package documentation -> https://github.com/go-chi/chi/tree/master/_examples/graceful
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal(":: Graceful shutdown timed out ... forcing exit! ::")
			}
		}()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		cancel() // <- it's not necessary!
		serverStopCtx()
	}()
	log.Println("---> RUN SERVER <---")

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("XXX <-- SERVER STOPPED! --> XXX")
	<-serverCtx.Done()
}
