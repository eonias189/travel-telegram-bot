package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Location struct {
	Order     int    `json:"-"`
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

func (l Location) MarshalBinary() (data []byte, err error) {
	return json.Marshal(l)
}

type LocationService struct {
	cli        *redis.Client
	tripPrefix string
}

func (ls LocationService) getTripKey(tripId int64) string {
	return fmt.Sprintf(`%v:%v`, ls.tripPrefix, tripId)
}

func (ls *LocationService) Get(tripId int64, order int) (Location, error) {
	data, err := ls.cli.JSONGet(context.TODO(), ls.getTripKey(tripId), fmt.Sprintf(`$.locations[%v]`, order-1)).Result()
	if err != nil {
		return Location{}, err
	}

	locs := []Location{}
	err = json.Unmarshal([]byte(data), &locs)

	if err != nil {
		return Location{}, err
	}

	if len(locs) == 0 {
		return Location{}, ErrNotFound
	}

	loc := locs[0]
	loc.Order = order

	return loc, nil
}

func (ls *LocationService) Set(tripId int64, order int, location Location) error {
	all, err := ls.GetAll(tripId)
	if err != nil {
		return err
	}
	all[order-1] = location
	return ls.cli.JSONSet(context.TODO(), ls.getTripKey(tripId), "$.locations", all).Err()
}

func (ls *LocationService) GetAll(tripId int64) ([]Location, error) {
	err := ls.cli.JSONArrLen(context.TODO(), ls.getTripKey(tripId), "$.locations").Err()
	if errors.Is(err, redis.Nil) {
		return []Location{}, nil
	}

	if err != nil {
		return []Location{}, err
	}

	resp, err := ls.cli.JSONGet(context.TODO(), ls.getTripKey(tripId), "$.locations[*]").Result()
	if err != nil {
		return []Location{}, err
	}

	locations := []Location{}
	err = json.Unmarshal([]byte(resp), &locations)
	if err != nil {
		return []Location{}, err
	}

	for i := 0; i < len(locations); i++ {
		locations[i].Order = i + 1
	}

	return locations, err
}

func (ls *LocationService) Add(tripId int64, location Location) error {
	key := ls.getTripKey(tripId)
	was, _ := ls.cli.JSONGet(context.TODO(), key, "$.locations").Result()
	if was == "[null]" {
		return ls.cli.JSONSet(context.TODO(), key, "$.locations", []Location{location}).Err()
	}
	return ls.cli.JSONArrAppend(context.TODO(), key, "$.locations", location).Err()
}

func NewLocationService(cli *redis.Client) *LocationService {
	return &LocationService{cli: cli, tripPrefix: "trips"}
}

func (ls *LocationService) Delete(tripId int64, order int) error {
	return ls.cli.JSONArrPop(context.TODO(), ls.getTripKey(tripId), "$.locations", order-1).Err()
}
