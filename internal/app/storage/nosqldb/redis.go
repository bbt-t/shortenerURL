package nosqldb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bbt-t/shortenerURL/internal/app/storage"
	"log"
	"time"

	"github.com/bbt-t/shortenerURL/configs"

	"github.com/go-redis/redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func NewRedisConnect() storage.DBRepo {
	/*
		Connect to Redis.
	*/
	cfg := configs.NewConfRedis()
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHOST, cfg.RedisPORT),
			Password: cfg.RedisPASS,
			DB:       0,
		}),
	}
}

func (r *redisClient) SaveURL(k, v string) error {
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

func (r *redisClient) GetURL(k string) (string, error) {
	/*
		Get value by key.
		param k: search key
		return: found value and error (or nil)
	*/
	ctx := context.Background()
	val, err := r.client.Get(ctx, k).Result()
	if errors.Is(err, redis.Nil) {
		return "", err
	}
	return val, nil
}

func (r *redisClient) Ping() error {
	ctx := context.Background()

	status := r.client.Ping(ctx)
	err := status.Err()

	if err != nil {
		log.Println(err)
	}
	return err
}
