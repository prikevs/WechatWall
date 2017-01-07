package libredis

import (
	"gopkg.in/redis.v5"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var (
	Addr     = "localhost:6379"
	Password = "" // no password set
	DB       = 0  // use default DB
)

func Init(cfg *RedisConfig) {
	Addr = cfg.Addr
	Password = cfg.Password
	DB = cfg.DB
}

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
