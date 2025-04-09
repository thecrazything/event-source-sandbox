package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var (
	instance *RedisClient
	once     sync.Once
)

func GetRedisClient() *RedisClient {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			DB:       0,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		instance = &RedisClient{
			client: rdb,
			ctx:    ctx,
		}
	})
	return instance
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
