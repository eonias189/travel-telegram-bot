package geoapi

import (
	"fmt"

	"github.com/golang/geo/s2"
	"github.com/serjvanilla/go-overpass"
)

func GetAttractions(p s2.LatLng, round int) ([]string, error) {
	c := overpass.New()

	around := fmt.Sprintf(`(around:%v,%v,%v)`, round, p.Lat, p.Lng)

	query := fmt.Sprintf(
		`
	[out:json];
(
  node["tourism"="attraction"]["name:ru"]%v;
  node["tourism"="museum"]["name:ru"]%v;
  node["attraction"]["name:ru"]%v;
);
out body;
	`, around, around, around)

	resp, err := c.Query(query)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for _, n := range resp.Nodes {
		res = append(res, n.Tags["name:ru"])
	}

	return res, nil

}
