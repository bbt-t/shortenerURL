package app

import (
	"log"

	"github.com/bbt-t/shortenerURL/internal/adapter/storage"
	"github.com/bbt-t/shortenerURL/internal/config"
	"github.com/bbt-t/shortenerURL/internal/controller/rest"
	"github.com/bbt-t/shortenerURL/internal/controller/rest/handler"
	"github.com/bbt-t/shortenerURL/internal/usecase"
)

func Run(cfg *configs.ServerCfg) {
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
		repo = storage.NewMapDBPlug()
	}
	if err := repo.PingDB(); err != nil {
		log.Println(err)
	}

	service := usecase.NewShortener(repo)
	h := handler.NewShortenerRoutes(service, cfg)
	server := rest.NewHTTPServer(cfg.ServerAddress, h.InitRoutes())

	log.Fatal(server.UP())
}
