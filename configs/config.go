package configs

import (
	"fmt"
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

func NewConfServ() *ServerCfg {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	return &ServerCfg{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseURL:       os.Getenv("BASE_URL"),
	}
}

type RedisConfig struct {
	RedisHOST string
	RedisPORT string
	RedisPASS string
}

func NewConfRedis() *RedisConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	return &RedisConfig{
		RedisHOST: os.Getenv("REDIS_HOST"),
		RedisPASS: os.Getenv("REDIS_PASS"),
		RedisPORT: os.Getenv("REDIS_PORT"),
	}
}

type SQLiteConfig struct {
	DBName string
}

func NewConfSQLite() *SQLiteConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	return &SQLiteConfig{
		DBName: os.Getenv("DB_NAME"),
	}
}

type PGConfig struct {
	DBUrl string
}

func NewConfPG() *PGConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	return &PGConfig{
		DBUrl: fmt.Sprintf(
			"host=%s dbname=%s user=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
		),
	}
}
