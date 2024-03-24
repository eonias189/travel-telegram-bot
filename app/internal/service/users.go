package service

import (
	"context"

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

func (us *UserService) AddTrip(userId, tripId int64) error {
	was, _ := us.cli.JSONGet(context.TODO(), us.getKey(userId), "$.trips").Result()
	if was == "[null]" {
		return us.cli.JSONSet(context.TODO(), us.getKey(userId), "$.trips", []int64{tripId}).Err()
	}
	return us.cli.JSONArrAppend(context.TODO(), us.getKey(userId), "$.trips", tripId).Err()
}

func (us *UserService) DeleteTrip(userId, tripId int64) error {
	tripsToDelete, err := us.cli.JSONArrIndex(context.TODO(), us.getKey(userId), "$.trips", tripId).Result()
	if err != nil {
		return err
	}

	if len(tripsToDelete) == 0 {
		return ErrNotFound
	}

	return us.cli.JSONArrPop(context.TODO(), us.getKey(userId), "$.trips", int(tripsToDelete[0])).Err()
}

func NewUserServive(cli *redis.Client) *UserService {
	return &UserService{JsonMixin: JsonMixin[User]{cli: cli, prefix: "users"}}
}
