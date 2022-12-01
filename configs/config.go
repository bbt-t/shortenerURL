package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerCfg struct {
	ServerAddress string
	BaseURL       string
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	sa := os.Getenv("SERVER_ADDRESS")
	bu := os.Getenv("BASE_URL")
	if sa == "" {
		sa = "127.0.0.1"
	}
	if bu == "" {
		bu = "http://127.0.0.1:8080"
	}
	return &ServerCfg{
		ServerAddress: sa,
		BaseURL:       bu,
	}
}

type RedisConfig struct {
	RedisHOST string
	RedisPORT string
	RedisPASS string
}

func NewConfRedis() *RedisConfig {
	/*
		Initialize a new Redis conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	rh := os.Getenv("REDIS_HOST")
	rp := os.Getenv("REDIS_PORT")
	if rh == "" {
		rh = "127.0.0.1"
	}
	if rp == "" {
		rp = "6379"
	}
	return &RedisConfig{
		RedisHOST: rh,
		RedisPASS: os.Getenv("REDIS_PASS"),
		RedisPORT: rp,
	}
}

type SQLiteConfig struct {
	DBName string
}

func NewConfSQLite() *SQLiteConfig {
	/*
		Initialize a new SQLite DB conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	var name string

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	if name = os.Getenv("DB_NAME"); name == "" {
		name = "./DefaultDBName.db"
	}
	return &SQLiteConfig{
		DBName: name,
	}
}

type PGConfig struct {
	DBUrl string
}

func NewConfPG() *PGConfig {
	//TODO Fully implement
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
