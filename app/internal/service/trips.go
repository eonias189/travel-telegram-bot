package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Trip struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Creator     int64      `json:"creator"`
	Members     []int64    `json:"members"`
	Locations   []Location `json:"locations"`
}

func initTrips(rdb *redis.Client) error {
	query := `
	FT.CREATE tripsIdx
	ON JSON
	PREFIX 1 trips:
	SCHEMA
		$.name as name TEXT
	`

	err := rdb.Do(context.TODO(), parseCommand(query)...).Err()
	if err != nil && err.Error() == "Index already exists" {
		return nil
	}

	return err
}

type TripService struct {
	JsonMixin[Trip]
}

func (ts *TripService) GetAll(ids []int64) ([]Trip, error) {
	res := make([]Trip, len(ids))

	for i, id := range ids {

		trip, err := ts.Get(id)
		if err != nil {
			return []Trip{}, err
		}

		res[i] = trip
	}

	return res, nil
}

func (ts *TripService) ExistsName(name string) bool {
	res, err := ts.cli.Do(context.TODO(), parseCommand(fmt.Sprintf(`FT.SEARCH tripsIdx @name:%v`, name))...).Result()

	if err != nil {
		return false
	}

	resMap, _ := res.(map[any]any)
	return resMap["total_results"] != int64(0)
}

func NewTripService(cli *redis.Client) *TripService {
	return &TripService{JsonMixin: JsonMixin[Trip]{cli: cli, prefix: "trips"}}
}
