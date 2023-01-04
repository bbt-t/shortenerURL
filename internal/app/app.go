package app

import (
	"github.com/bbt-t/shortenerURL/internal/adapter/storage"
	"github.com/bbt-t/shortenerURL/internal/config"
	"github.com/bbt-t/shortenerURL/internal/controller/rest"
	"github.com/bbt-t/shortenerURL/internal/controller/rest/handler"
	"github.com/bbt-t/shortenerURL/internal/usecase"
	"log"
)

func Run(cfg *configs.ServerCfg) {
	/*
		Создание используемых объектов через конструкторы для слоёв
		и graceful shutdown.
	*/
	var repo storage.DatabaseRepository

	//if fp := cfg.FilePath; fp != "" {
	//	repo = storage.NewFileDB(fp)
	//}
	//if dsn := cfg.DBConnectURL; dsn != "" {
	//	repo = storage.NewSQLDatabase(dsn)
	//}
	//// Default in memory storage:
	//if err := repo.PingDB(); err != nil {
	//	repo = storage.NewMapDBPlug()
	//}

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
	server := rest.NewHttpServer(cfg.ServerAddress, h.InitRoutes())

	//go func() {
	//	if err := server.UP(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("server :: %s :: DIED, reason --> %v", cfg.BaseURL, err)
	//	}
	//}()
	log.Fatal(server.UP())
}
