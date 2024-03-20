package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCash struct {
	conn      *redis.Conn
	expirTime time.Duration
}

func (rc *RedisCash) Set(key string, value string) error {
	return rc.conn.Set(context.TODO(), key, value, rc.expirTime).Err()
}

func (rc *RedisCash) Get(key string) (string, error) {
	return rc.conn.Get(context.TODO(), key).Result()
}

func NewRedisCash(conn *redis.Conn, expirationTime time.Duration) *RedisCash {
	return &RedisCash{conn: conn, expirTime: expirationTime}
}
