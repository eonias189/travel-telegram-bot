package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect(address, user, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: user,
		Password: password,
		DB:       db,
	})
	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
