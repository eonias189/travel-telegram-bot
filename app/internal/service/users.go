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

func NewUserServive(cli *redis.Client) *UserService {
	return &UserService{JsonMixin: JsonMixin[User]{cli: cli, prefix: "users"}}
}
