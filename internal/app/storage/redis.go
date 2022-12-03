package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bbt-t/shortenerURL/configs"

	"github.com/go-redis/redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisConnect() *RedisClient {
	/*
		Connect to Redis.
	*/
	cfg := configs.NewConfRedis()
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHOST, cfg.RedisPORT),
			Password: cfg.RedisPASS,
			DB:       0,
		}),
	}
}

func (r RedisClient) SaveURL(k, v string) error {
	/*
		Write key - value to Redis.
	*/
	ctx := context.Background()
	err := r.client.Set(ctx, k, v, 20*time.Second).Err()
	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	return err
}

func (r RedisClient) GetURL(k string) (string, error) {
	/*
		Get value by key.
		param k: search key
		return: found value and error (or nil)
	*/
	ctx := context.Background()
	val, err := r.client.Get(ctx, k).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}
