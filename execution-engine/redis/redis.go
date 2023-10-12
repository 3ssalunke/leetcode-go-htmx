package redis

import (
	"time"

	"github.com/3ssalunke/leetcode-clone-exen/util"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(config util.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	return &RedisClient{Client: client}
}

func (client *RedisClient) Ping() (string, error) {
	return client.Client.Ping().Result()
}

func (client *RedisClient) SetValue(key string, value interface{}) error {
	return client.Client.Set(key, value, time.Hour).Err()
}
