package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-12543.c273.us-east-1-2.ec2.cloud.redislabs.com:12543",
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return rdb
}

func SetKey(ctx *context.Context, rdb *redis.Client, key string, value string, ttl int) {
	rdb.Set(*ctx, key, value, 0)
}

func GetLongURL(ctx *context.Context, rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := rdb.Get(*ctx, shortURL).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}
	return longURL, nil
}
