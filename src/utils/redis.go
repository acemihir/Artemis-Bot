package utils

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

var Cache *Redis

func SetupCache() error {
	Cache = &Redis{
		Context: context.Background(),
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	// Perform little test
	_, ex := Cache.Client.Ping(Cache.Context).Result()
	if ex != nil {
		log.Fatalln("[ERROR] Could not ping redis.")
	}

	return ex
}

func (at *Redis) SetCache(key, val string) {
	ex := at.Client.Set(at.Context, key, val, 0).Err()
	if ex != nil {
		log.Fatalln("[ERROR] Could not set in redis.")
	}
}

func (at *Redis) ExistsCache(key string) int64 {
	res, ex := at.Client.Exists(at.Context, key).Result()
	if ex != nil {
		log.Fatalln("[ERROR] Could not check for existance in redis.")
	}
	return res
}

func (at *Redis) GetCache(key string) string {
	res, ex := at.Client.Get(at.Context, key).Result()
	if ex != nil {
		log.Fatalln("[ERROR] Could not get from redis.")
	}
	return res
}

func (at *Redis) DelCache(key string) {
	ex := at.Client.Del(at.Context, key).Err()
	if ex != nil {
		log.Fatalln("[ERROR] Could not delete from redis.")
	}
}
