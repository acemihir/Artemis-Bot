package utils

import (
	"context"
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
	_, ex := Cache.Client.Ping(Cache.Context).Result()
	if ex != nil {
		Cout("[ERROR] Redis ping failed: %v", Red, ex)
	}
}

func (at *Redis) SetCache(key, val string) {
	ex := at.Client.Set(at.Context, key, val, at.Expiry).Err()
	if ex != nil {
		Cout("[ERROR] Set in redis failed: %v", Red, ex)
	}
}

func (at *Redis) ExistsCache(key string) int64 {
	res, ex := at.Client.Exists(at.Context, key).Result()
	if ex != nil {
		Cout("[ERROR] Existance check in redis failed: %v", Red, ex)
	}
	return res
}

func (at *Redis) GetCache(key string) string {
	res, ex := at.Client.Get(at.Context, key).Result()
	if ex != nil {
		Cout("[ERROR] Get from redis failed: %v", Red, ex)
	}
	return res
}

func (at *Redis) DelCache(key string) {
	ex := at.Client.Del(at.Context, key).Err()
	if ex != nil {
		Cout("[ERROR] Delete from redis failed: %v", Red, ex)
	}
}
