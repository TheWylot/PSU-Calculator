package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Vigiatonet/PSU-Calculator/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func LoadRedis(cfg *config.Config) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		PoolSize:     cfg.Redis.PoolSize,
	})

	_, err := redisClient.Ping(context.Background()).Result() // first result is PONG
	if err != nil {
		return err
	}
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	err := redisClient.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Set[T any](key string, value T, ttl time.Duration, c *redis.Client) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.Set(context.Background(), key, v, ttl).Result() // first result is OK
	if err != nil {
		return err
	}
	return nil

}

func Get[T any](key string, c *redis.Client) (*T, error) {
	var result T = *new(T)
	resultByte, err := c.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resultByte), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
