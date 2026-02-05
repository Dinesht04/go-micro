package data

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// type RedisClient struct {
// 	Client *redis.NewClient
// }

func NewRedisClient(ctx context.Context, logger *slog.Logger) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return rdb, err
	}

	return rdb, nil

}
