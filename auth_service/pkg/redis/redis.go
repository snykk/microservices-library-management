package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	Set(key string, value interface{}) error
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error

	Ping() error
}

type redisCache struct {
	expires time.Duration
	client  *redis.Client
}

func NewRedisCache(host string, db int, password string, expires time.Duration) RedisCache {
	return &redisCache{
		expires: expires,
		client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       db,
		}),
	}
}

func (cache *redisCache) Ping() error {
	return cache.client.Ping(context.Background()).Err()
}

func (cache *redisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cache.client.Set(cache.client.Context(), key, json, cache.expires).Err()
}

func (cache *redisCache) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cache.client.Set(cache.client.Context(), key, jsonData, expiration).Err()
}

func (cache *redisCache) Get(key string) (email string, err error) {
	val, err := cache.client.Get(cache.client.Context(), key).Result()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(val), &email)
	return email, err
}

func (cache *redisCache) Del(key string) error {
	return cache.client.Del(cache.client.Context(), key).Err()
}
