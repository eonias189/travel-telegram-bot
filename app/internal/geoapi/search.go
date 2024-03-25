package geoapi

import (
	"errors"
	"fmt"

	"github.com/golang/geo/s2"
	"github.com/yuriizinets/go-nominatim"
)

var (
	ErrNotFound = errors.New("NotFound")
)

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

func GetCoords(location string) (s2.LatLng, error) {
	n := nominatim.Nominatim{}
	res, err := n.Search(nominatim.SearchParameters{
		Query: location,
	})

	if err != nil {
		return s2.LatLng{}, err
	}

	if len(res) == 0 {
		return s2.LatLng{}, ErrNotFound
	}

	return s2.LatLngFromDegrees(res[0].Lat, res[0].Lng), nil
}

type Address struct {
	Country string
	City    string
}

func GetAddress(p s2.LatLng) (nominatim.SearchAddress, error) {
	n := nominatim.Nominatim{}
	res, err := n.Search(nominatim.SearchParameters{
		Query:          fmt.Sprintf(`%v,%v`, p.Lat, p.Lng),
		IncludeAddress: true,
	})

	if err != nil {
		return nominatim.SearchAddress{}, err
	}

	if len(res) == 0 {
		return nominatim.SearchAddress{}, ErrNotFound
	}
	return res[0].Address, nil
}
