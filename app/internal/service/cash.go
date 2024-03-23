package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	errlogger "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/loggers/errLogger"
	"github.com/redis/go-redis/v9"
)

type RedisCash struct {
	conn      *redis.Conn
	expirTime time.Duration
	prefix    string
}

func (rc RedisCash) getKey(key string) string {
	return fmt.Sprintf("%v:%v", rc.prefix, key)
}

func (rc *RedisCash) Set(key string, value string) {
	err := rc.conn.Set(context.TODO(), rc.getKey(key), value, rc.expirTime).Err()
	if err != nil {
		errlogger.New().Error("unable to set context", slog.String("err", err.Error()))
	}
}

func (rc *RedisCash) Get(key string) (string, bool) {
	resp := rc.conn.Get(context.TODO(), rc.getKey(key))
	if err := resp.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false
		}
		errlogger.New().Error("unable to get context", slog.String("err", err.Error()))
		return "", false
	}
	return resp.Val(), true
}

type CashOptions struct {
	Prefix         string
	ExpirationTime time.Duration
}

func NewRedisCash(conn *redis.Conn, opts CashOptions) *RedisCash {
	return &RedisCash{conn: conn, prefix: opts.Prefix, expirTime: opts.ExpirationTime}
}
