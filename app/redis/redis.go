package redis

import (
	"github.com/3ssalunke/leetcode-clone-app/util"
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

func (client *RedisClient) GetValue(key string) (string, error) {
	return client.Client.Get(key).Result()
}
