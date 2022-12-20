package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bbt-t/shortenerURL/configs"

	"github.com/go-redis/redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func NewRedisConnect() DBRepo {
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
		Write key - value (strings) to Redis.
	*/
	ctx := context.Background()
	err := r.client.Set(ctx, k, v, 20*time.Second).Err()

	if err != nil {
		log.Printf("ERROR : %s", err)
	}
	return err
}

func (r *redisClient) SaveUser(k string, value interface{}) error {
	/*
		Save user info in gob-representation.
		param k: key in string
		param value: struct obj
	*/
	var buf bytes.Buffer
	ctx := context.Background()
	encoder := gob.NewEncoder(&buf)
	id := "user-db"

	if err := encoder.Encode(value); err != nil {
		log.Println(err)
		return err
	}
	if err := r.client.HSet(ctx, id, k, buf.Bytes()).Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
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
		log.Printf("key %v Not found", k)
		return "", err
	}
	return val, nil
}

func (r *redisClient) GetUser(k string, value interface{}) error {
	/*
		Get value by key.
	*/
	ctx := context.Background()
	id := "user-db"

	buf, err := r.client.HGet(ctx, id, k).Bytes()
	if err != nil {
		log.Println(err)
		return err
	}
	errDec := gob.NewDecoder(bytes.NewReader(buf)).Decode(value)
	if errDec != nil {
		log.Println(err)
		return errDec
	}
	return nil
}

func (r *redisClient) Ping() error {
	ctx := context.Background()
	status := r.client.Ping(ctx)
	if err := status.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
