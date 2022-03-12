package utils

import (
	"context"
	"os"
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
		os.Exit(1)
	}
}

func (at *Redis) SetCache(key, val string) error {
	return at.Client.Set(at.Context, key, val, at.Expiry).Err()
}

func (at *Redis) ExistsCache(key string) (int64, error) {
	return at.Client.Exists(at.Context, key).Result()
}

func (at *Redis) GetCache(key string) (string, error) {
	return at.Client.Get(at.Context, key).Result()
}

func (at *Redis) DelCache(key string) error {
	return at.Client.Del(at.Context, key).Err()
}
