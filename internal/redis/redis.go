package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, target interface{}) error
	Delete(ctx context.Context, key string) error
}

type redisService struct {
	client *redis.Client
}

func NewRedisService(addr string, password string, db int) RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisService{
		client: rdb,
	}
}

func (r *redisService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis.Set: failed to marshal: %w", err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisService) Get(ctx context.Context, key string, target interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key does not exist")
	} else if err != nil {
		return fmt.Errorf("redis.Get: %w", err)
	}

	return json.Unmarshal([]byte(val), target)
}

func (r *redisService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
