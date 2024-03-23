package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Age      int    `json:"age"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
}

func initUsers(rdb *redis.Client) error {
	query := `
	FT.CREATE usersIdx
	ON JSON
	PREFIX 1 users:
	SCHEMA
		$.age as age NUMERIC
		$.location as location TEXT
		$.bio as bio TEXT
	`

	cmd := rdb.Do(context.TODO(), parseCommand(query)...)
	err := cmd.Err()
	if err != nil && err.Error() == "Index already exists" {
		return nil
	}
	return err
}

type UserService struct {
	conn *redis.Conn
}

func (us *UserService) Get(id int64) (User, error) {
	data, err := us.conn.JSONGet(context.TODO(), fmt.Sprintf("users:%v", id)).Result()
	if err != nil {
		return User{}, err
	}

	if data == "" {
		return User{}, ErrNotFound
	}

	user := User{}
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (us *UserService) Set(id int64, user User) error {
	return us.conn.JSONSet(context.TODO(), fmt.Sprintf("users:%v", id), "$", user).Err()
}

func NewUserServive(conn *redis.Conn) *UserService {
	return &UserService{conn: conn}
}
