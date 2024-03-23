package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	saveInterval = time.Second * 5
)

func Connect(ctx context.Context, address, user, password string, db int) (*redis.Client, error) {
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rdb.Save(context.TODO())
				time.Sleep(saveInterval)
			}
		}
	}()

	return rdb, nil
}
