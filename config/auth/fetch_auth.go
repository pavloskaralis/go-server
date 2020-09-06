package auth

import (
	"github.com/go-redis/redis/v7"
	"os"
)

func FetchAuth(authD *AccessDetails) (string, error) {
	//connect to redis
	var  client *redis.Client
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	userID, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
	   return "", err
	}
	return userID, nil
  }