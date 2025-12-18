package svc

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(c *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// 测试Redis连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
