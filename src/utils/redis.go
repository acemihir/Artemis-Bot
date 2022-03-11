package utils

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
	Expiry  time.Duration
}

var Cache *Redis

func SetupCache(expiryTime time.Duration) {
	Cache = &Redis{
		Context: context.Background(),
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
		Expiry: expiryTime,
	}

	// Perform little test
	_, pEx := Cache.Client.Ping(Cache.Context).Result()
	if pEx != nil {
		log.Fatalf("[ERROR] Redis ping failed: %s", pEx)
	}
}

func (at *Redis) SetCache(key, val string) {
	ex := at.Client.Set(at.Context, key, val, at.Expiry).Err()
	if ex != nil {
		log.Fatalf("[ERROR] Set in redis failed: %s", ex)
	}
}

func (at *Redis) ExistsCache(key string) int64 {
	res, ex := at.Client.Exists(at.Context, key).Result()
	if ex != nil {
		log.Fatalf("[ERROR] Existance check in redis failed: %s", ex)
	}
	return res
}

func (at *Redis) GetCache(key string) string {
	res, ex := at.Client.Get(at.Context, key).Result()
	if ex != nil {
		log.Fatalf("[ERROR] Get from redis failed: %s", ex)
	}
	return res
}

func (at *Redis) DelCache(key string) {
	ex := at.Client.Del(at.Context, key).Err()
	if ex != nil {
		log.Fatalf("[ERROR] Delete from redis failed: %s", ex)
	}
}
