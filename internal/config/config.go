package configs

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerCfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"    envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL"          envDefault:"http://127.0.0.1:8080"`
	FilePath      string `env:"FILE_STORAGE_PATH"` //envDefault:"FILE_OBJ.gob"`
	DBConnectURL  string `env:"DATABASE_DSN"`      //envDefault:"host=localhost port=5432 user=postgres password=$apr1$dISdUBfu$NCBQX/q3R2WUV1JppxP8l0 dbname=postgres sslmode=disable"`
	DBused        string
	//SecretKey     string
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf.
		flag -> env, env-variables take precedence.
	*/
	cfg := ServerCfg{}

	flag.StringVar(&cfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "", "base url")
	flag.StringVar(&cfg.FilePath, "f", "", "file path")
	flag.StringVar(&cfg.DBConnectURL, "d", "", "db url (only for pg")
	//flag.StringVar(&cfg.SecretKey, "k", "", "secret key to sign uid cookies")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Database selection by priority:
	if cfg.FilePath != "" {
		cfg.DBused = "file"
	}
	if cfg.DBConnectURL != "" {
		cfg.DBused = "pg"
	}

	return &cfg
}
