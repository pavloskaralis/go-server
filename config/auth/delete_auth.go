
package auth

import (
	"os"
	"github.com/go-redis/redis/v7"
)

func DeleteAuth(givenUuid string) (int64,error) {
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

	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
	   return 0, err
	}
	return deleted, nil
}
