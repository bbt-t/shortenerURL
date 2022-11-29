package storage

import (
	"context"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/go-redis/redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func RedisClientConnect() *redis.Client {
	/*
		Connect to Redis.
	*/
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.NewConfRedis().RedisHOST, configs.NewConfRedis().RedisPORT),
		Password: configs.NewConfRedis().RedisPASS,
		DB:       0,
	})
}

func (r RedisClient) SaveURL(k, v string) error {
	/*
		Write key - value to Redis.
	*/
	ctx := context.Background()
	err := r.client.Set(ctx, k, v, 20*time.Second).Err()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
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
