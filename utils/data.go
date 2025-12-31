package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
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

func NewMongoClient() (*mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SaveURLToDB(ctx context.Context, client *mongo.Client, code, longURL string) error {
	collection := client.Database("ziplink").Collection("main")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       code,
		"long_url":   longURL,
		"created_at": time.Now(),
	})
	return err
}

func GetURLFromDB(ctx context.Context, client *mongo.Client, code string) (string, error) {
	collection := client.Database("ziplink").Collection("main")
	var result struct {
		LongURL string `bson:"long_url"`
	}
	err := collection.FindOne(ctx, bson.M{"code": code}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.LongURL, nil
}
