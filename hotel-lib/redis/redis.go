package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	db *redis.Client
}

func NewClient() (*RedisClient, error) {
	log := logger.New()
	db := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		log.Logger.Info(fmt.Sprintf("failed to connect to redis server: %s\n", err.Error()))
		return nil, err
	}
	return &RedisClient{
		db: db,
	}, nil
}

func (r *RedisClient) SaveJSON(key string, value []byte) error {
	err := r.db.Set(context.Background(), key, value, 30*time.Second)

	if err != nil {
		return err.Err()
	}
	return nil
}

func (r *RedisClient) GetData(key string) (bool, string) {
	val, err := r.db.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, ""
		} 
	}
	return true, val
}