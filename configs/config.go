package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerCfg struct {
	ServerAddress string
	BaseURL       string
	RedisPASS     string
	RedisPORT     string
}

func NewConfig() *ServerCfg {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	return &ServerCfg{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseURL:       os.Getenv("BASE_URL"),
		RedisPASS:     os.Getenv("REDIS_PASS"),
		RedisPORT:     os.Getenv("REDIS_PORT"),
	}
}
