package configs

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerCfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"`    // envDefault:"127.0.0.1:8080"
	BaseURL       string `env:"BASE_URL"`          //envDefault:"http://127.0.0.1:8080"
	FilePath      string `env:"FILE_STORAGE_PATH"` //envDefault:"FILE_OBJ.gob"
	UseDB         string
}

func NewConfServ() *ServerCfg {
	/*
		Initialize a new conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	var cfg ServerCfg

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	flag.String("db", "none", `
	"select database: no flag - use map, 'sqlite' - use SQLite, 'pg' - Postgresql, 'redis' - Redis"
	`)
	flag.String("a", "", "server address")
	flag.String("b", "", "base url")
	flag.String("f", "", "file path")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = viper.GetString("a")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = viper.GetString("b")
	}
	if cfg.FilePath == "" {
		cfg.FilePath = viper.GetString("f")
	}
	if cfg.UseDB == "" {
		cfg.UseDB = viper.GetString("db")
	}

	return &cfg
}

type RedisConfig struct {
	RedisHOST string `toml:"redis_host"`
	RedisPORT string `toml:"redis_port"`
	RedisPASS string `env:"REDIS_PASS,file" envDefault:""`
}

func NewConfRedis() *RedisConfig {
	/*
		Initialize a new Redis conf. Values are taken from .env file.
		If .env file does not exist or the required value does not exist,
		then default values are substituted.
	*/
	var db RedisConfig
	_, err := toml.DecodeFile("./configs/server_conf.toml", &db)
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	if err := env.Parse(&db); err != nil {
		fmt.Printf("%+v\n", err)
	}

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
	_, err := toml.DecodeFile("./configs/server_conf.toml", &db)
	if err != nil {
		db.DBName = "./data/DefaultName.db"
		log.Printf("ERROR : %s", err)
	}
	return &db
}

type PGConfig struct {
	DBUrl string
}

func NewConfPG(param ...bool /* optional args */) *PGConfig {
	/*
		return: url-param for connect to PG DB.
	*/
	var pgCfg pgMakeConfig

	if err := env.Parse(&pgCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &PGConfig{
		DBUrl: pgCfg.makeURL(),
	}
}

type pgMakeConfig struct {
	host     string `env:"DB_HOST"`
	dbname   string `env:"DB_NAME"`
	user     string `env:"DB_USER"`
	password string `env:"DB_PASSWORD"`
	sslMode  string `env:"SSL_MODE"`
}

func (p *pgMakeConfig) makeURL() string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		p.host,
		p.dbname,
		p.user,
		p.password,
	)
}
