package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerCfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"    envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL"          envDefault:"http://127.0.0.1:8080"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBConnectURL  string `env:"DATABASE_DSN" envDefault:"host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"`
	DBused        string
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new config.
	*/
	var cfg ServerCfg

	flag.StringVar(&cfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "", "base url")
	flag.StringVar(&cfg.FilePath, "f", "", "file path")
	flag.StringVar(&cfg.DBConnectURL, "d", "", "postgres DSN (url)")

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	flag.Parse()

	// Database selection by priority:
	if cfg.FilePath != "" {
		cfg.DBused = "file"
	}
	if cfg.DBConnectURL != "" {
		cfg.DBused = "pg"
	}

	return &cfg
}
