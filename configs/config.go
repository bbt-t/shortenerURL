package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type ServerCfg struct {
	ServerAddress string `toml:"address"`
	Port          string `toml:"port"`
	BaseURL       string `toml:"baseurl"`
	LogLevel      string `toml:"loglevel"`
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	var conf ServerCfg
	_, err := toml.DecodeFile("./configs/serverconf.toml", &conf)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	return &conf
}

type RedisConfig struct {
	RedisHOST string `toml:"redis_host"`
	RedisPORT string `toml:"redis_port"`
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
	var db RedisConfig
	_, err := toml.DecodeFile("./configs/serverconf.toml", &db)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	db.RedisPASS = os.Getenv("REDIS_PASS")
	return &db
}

type SQLiteConfig struct {
	DBName string `toml:"db_file_name"`
}

func NewConfSQLite() *SQLiteConfig {
	/*
		Initialize a new SQLite DB conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	var db SQLiteConfig
	_, err := toml.DecodeFile("./configs/serverconf.toml", &db)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	return &db
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
