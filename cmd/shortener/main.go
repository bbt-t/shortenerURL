package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/go-redis/redis/v9"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	cfg := configs.NewConfig()

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:8080", cfg.ServerAddress),
		MaxHeaderBytes:    0,
		ReadHeaderTimeout: 1 * time.Second,
	}
	http.HandleFunc("/", RedirectToOriginalURL)
	log.Fatal(server.ListenAndServe())
}

func RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p := strings.Split(r.URL.Path, "/")[1]
		rUrl, err := GetVarRedis(RedisClientConnect(), p)

		if err != nil {
			log.Printf("ERROR : %s", err)
			http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)
		}
		w.Header().Set("Location", rUrl)
		w.WriteHeader(307)
		http.Redirect(w, r, rUrl, http.StatusTemporaryRedirect)

	case http.MethodPost:
		defer r.Body.Close()
		var value CreateShortURLRequest

		payload, errReadBody := io.ReadAll(r.Body)
		if errReadBody != nil {
			log.Printf("ERROR : %s", errReadBody)
		}

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("ERROR: %s", err)
		}
		fmt.Println(value.URL)

		toHashVar := fmt.Sprintf("%d", HashShortening([]byte(value.URL)))
		SetVarRedis(RedisClientConnect(), toHashVar, value.URL)

		type Resp struct {
			Result string `json:"result"`
		}
		resp := Resp{
			Result: configs.NewConfig().BaseURL + "/" + toHashVar,
		}
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(res); err != nil {
			log.Printf("ERROR : %s", err)
		}
	}
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

func HashShortening(s []byte) uint32 {
	hash := fnv.New32a()
	if _, err := hash.Write(s); err != nil {
		log.Fatalf("ERROR : %s", err)
	}
	return hash.Sum32()
}

func RedisClientConnect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.NewConfig().ServerAddress, configs.NewConfig().RedisPORT),
		Password: configs.NewConfig().RedisPASS,
		DB:       0,
	})
}

var ctx = context.Background()

func SetVarRedis(rdb *redis.Client, k, v string) {
	err := rdb.Set(ctx, k, v, 20*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func GetVarRedis(rdb *redis.Client, k string) (string, error) {
	val, err := rdb.Get(ctx, k).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}
