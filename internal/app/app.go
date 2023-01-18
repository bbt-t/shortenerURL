package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbt-t/shortenerURL/internal/adapter/storage"
	"github.com/bbt-t/shortenerURL/internal/config"
	"github.com/bbt-t/shortenerURL/internal/controller/rest"
	"github.com/bbt-t/shortenerURL/internal/controller/rest/handler"
	"github.com/bbt-t/shortenerURL/internal/usecase"
)

func Run(cfg *config.ServerCfg) {
	/*
		Creating usable objects via constructors for layers and start app.
	*/
	var repo storage.DatabaseRepository

	switch cfg.DBused {
	case "pg":
		repo = storage.NewSQLDatabase(cfg.DBConnectURL)
	case "file":
		repo = storage.NewFileDB(cfg.FilePath)
	default:
		repo = storage.NewMapDB()
	}
	if err := repo.PingDB(); err != nil {
		log.Fatal(err)
	}

	service := usecase.NewShortener(repo)
	h := handler.NewShortenerRoutes(service, cfg)
	server := rest.NewHTTPServer(cfg.ServerAddress, h.InitRoutes())

	go func() {
		log.Println(server.UP())
	}()
	// Graceful shutdown:
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-gracefulStop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Printf("! Error shutting down server: !\n%v", err)
	} else {
		log.Println("! SERVER STOPPED !")
	}
}
