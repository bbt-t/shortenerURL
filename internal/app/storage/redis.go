package storage

import (
	"context"
	"fmt"
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/go-redis/redis/v9"
	"time"
)

func RedisClientConnect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.NewConfig().ServerAddress, configs.NewConfig().RedisPORT),
		Password: configs.NewConfig().RedisPASS,
		DB:       0,
	})
}

var ctx = context.Background()

func SaveNewUrlRedis(rdb *redis.Client, k, v string) {
	err := rdb.Set(ctx, k, v, 20*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func PullOutUrlRedis(rdb *redis.Client, k string) (string, error) {
	val, err := rdb.Get(ctx, k).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}
