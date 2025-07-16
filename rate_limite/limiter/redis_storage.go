package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(addr, password string, db int) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisStorage) Increment(key string, expire time.Duration) (int64, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		_ = r.client.Expire(r.ctx, key, expire).Err()
	}
	return count, nil
}

func (r *RedisStorage) Expire(key string, expire time.Duration) error {
	return r.client.Expire(r.ctx, key, expire).Err()
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
