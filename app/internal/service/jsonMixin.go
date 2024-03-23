package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type JsonMixin[T any] struct {
	prefix string
	conn   *redis.Conn
}

func (jm JsonMixin[T]) getKey(key int64) string {
	return fmt.Sprintf("%v:%v", jm.prefix, key)
}

func (jm *JsonMixin[T]) Get(key int64) (T, error) {
	var res T

	resp := jm.conn.JSONGet(context.TODO(), jm.getKey(key))
	if err := resp.Err(); err != nil {
		return res, err
	}

	data, err := resp.Result()
	if err != nil {
		return res, err
	}

	if data == "" {
		return res, ErrNotFound
	}

	err = json.Unmarshal([]byte(data), &res)
	return res, err
}

func (jm *JsonMixin[T]) Set(key int64, item T) error {
	return jm.conn.JSONSet(context.TODO(), jm.getKey(key), "$", item).Err()
}
