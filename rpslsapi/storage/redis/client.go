package redis

import (
	"github.com/go-redis/redis/v8"
	"rpsls/rpslsapi"
)

type Client struct {
	*redis.Client
}

func NewClient() Client {
	return Client{redis.NewClient(&redis.Options{
		Addr:     rpslsapi.Config.Redis.Addr,
		Password: rpslsapi.Config.Redis.Password,
		DB:       rpslsapi.Config.Redis.DB,
	})}
}
