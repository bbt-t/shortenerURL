package config

import (
	"fmt"
	flag "github.com/spf13/pflag"

	"github.com/caarlos0/env/v6"
)

type ServerCfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"    envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL"          envDefault:"http://127.0.0.1:8080"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBConnectURL  string `env:"DATABASE_DSN"`
	DBused        string
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf.
		flag -> env, env-variables take precedence.
	*/
	var cfg ServerCfg

	flag.StringVarP(&cfg.ServerAddress, "address", "a", "", "server address")
	flag.StringVarP(&cfg.BaseURL, "base", "b", "", "base url")
	flag.StringVarP(&cfg.FilePath, "file", "f", "", "file path")

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	flag.Parse()

	return &ServerCfg{
		ServerAddress: cfg.ServerAddress,
		BaseURL:       cfg.BaseURL,
		FilePath:      cfg.FilePath,
	}
}
