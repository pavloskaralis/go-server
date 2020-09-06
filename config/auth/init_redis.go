package auth

import (
	"github.com/go-redis/redis/v7"
	"os"
)

var Redis *redis.Client

func InitRedis() ( error ) {
	//connect to redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: dsn, 
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	Redis = client
	return nil
}