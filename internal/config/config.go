package config

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
	SecretKey     string
}

type FlagConfig struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	SecretKey       string
}

var flagCfg = FlagConfig{}

func init() {
	flag.StringVar(&flagCfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&flagCfg.BaseURL, "b", "", "base url")
	flag.StringVar(&flagCfg.FileStoragePath, "f", "", "file path")
	flag.StringVar(&flagCfg.SecretKey, "k", "", "secret key to sign uid cookies")
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf.
		flag -> env, env-variables take precedence.
	*/
	var cfg ServerCfg
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	cfg.UpdateFromFlags()

	return &cfg
}

func (cfg *ServerCfg) UpdateFromFlags() {
	if flagCfg.BaseURL != "" {
		cfg.BaseURL = flagCfg.BaseURL
	}
	if flagCfg.ServerAddress != "" {
		cfg.ServerAddress = flagCfg.ServerAddress
	}
	if flagCfg.SecretKey != "" {
		cfg.SecretKey = flagCfg.SecretKey
	}
}
