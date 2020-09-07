package auth

import (
	"github.com/go-redis/redis/v7"
)

var Redis *redis.Client

//called from main
func InitRedis() ( error ) {
	//connect to redis
	dsn := "localhost:6379"
	client := redis.NewClient(&redis.Options{
		Addr: dsn, 
	})
	//check connection
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	//set global var
	Redis = client
	return nil
}