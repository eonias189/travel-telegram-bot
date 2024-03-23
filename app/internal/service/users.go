package service

import (
	"github.com/redis/go-redis/v9"
)

type User struct {
	Age      int     `json:"age"`
	Location string  `json:"location"`
	Bio      string  `json:"bio"`
	Trips    []int64 `json:"trips"`
}

type UserService struct {
	JsonMixin[User]
}

func NewUserServive(conn *redis.Conn) *UserService {
	return &UserService{JsonMixin: JsonMixin[User]{conn: conn, prefix: "users:"}}
}
