package geoapi

import (
	"errors"

	"github.com/yuriizinets/go-nominatim"
)

var (
	ErrNotFound = errors.New("NotFound")
)

type Coords struct {
	Lat float64
	Lng float64
}

func CheckLocation(location string) bool {
	n := nominatim.Nominatim{}
	res, err := n.Search(nominatim.SearchParameters{
		Query: location,
	})

	if err != nil {
		return false
	}

	return len(res) > 0
}

func GetCoords(location string) (Coords, error) {
	n := nominatim.Nominatim{}
	res, err := n.Search(nominatim.SearchParameters{
		Query: location,
	})

	if err != nil {
		return Coords{}, err
	}

	if len(res) == 0 {
		return Coords{}, ErrNotFound
	}

	return Coords{Lat: res[0].Lat, Lng: res[0].Lng}, nil
}
