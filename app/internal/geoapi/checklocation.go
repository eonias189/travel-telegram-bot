package geoapi

import "github.com/yuriizinets/go-nominatim"

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
