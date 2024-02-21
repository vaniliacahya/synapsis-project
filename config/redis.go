package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

func InitRedisConfig() *redis.Client {
	redisDb, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 0, 16)
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       int(redisDb),                // use default DB
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Print("error redis")
		panic(err)
	}
	return rdb
}
