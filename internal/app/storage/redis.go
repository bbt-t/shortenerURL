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
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.NewConfRedis().RedisHOST, configs.NewConfRedis().RedisPORT),
		Password: configs.NewConfRedis().RedisPASS,
		DB:       0,
	})
}

var ctx = context.Background()

func (r RedisClient) SaveURL(k, v string) error {
	err := r.client.Set(ctx, k, v, 20*time.Second).Err()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
	}
	return err
}

func (r RedisClient) GetURL(k string) (string, error) {
	val, err := r.client.Get(ctx, k).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}
