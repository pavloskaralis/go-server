package auth

import (
	"github.com/go-redis/redis/v7"
	"time"
	"os"
)

func CreateAuth(userid string, td *TokenDetails) error {
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

	//unix to UTC
    at := time.Unix(td.AtExpires, 0) 
    rt := time.Unix(td.RtExpires, 0)
    now := time.Now()
	//store tokens
    errAccess := client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
    if errAccess != nil {
        return errAccess
    }
    errRefresh := client.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
    if errRefresh != nil {
        return errRefresh
    }
    return nil
}