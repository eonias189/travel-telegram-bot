package geoapi

import (
	"fmt"

	"github.com/golang/geo/s2"
	"github.com/serjvanilla/go-overpass"
)

type SearchResp struct {
	Name string
	P    s2.LatLng
}

func GetAttractions(p s2.LatLng, around int, limit int) ([]SearchResp, error) {
	c := overpass.New()

	ar := fmt.Sprintf(`(around:%v,%v,%v)`, around, p.Lat, p.Lng)

	query := fmt.Sprintf(
		`
	[out:json];
(
  node["tourism"="attraction"]["name:ru"]%v;
  node["tourism"="museum"]["name:ru"]%v;
  node["attraction"]["name:ru"]%v;
);
out body;
	`, ar, ar, ar)

	resp, err := c.Query(query)
	if err != nil {
		return []SearchResp{}, err
	}

	res := []SearchResp{}
	var count int
	for _, n := range resp.Nodes {
		if count == limit {
			break
		}
		res = append(res, SearchResp{Name: n.Tags["name:ru"], P: s2.LatLngFromDegrees(n.Lat, n.Lon)})
		count++
	}

	return res, nil

}

func GetHotels(p s2.LatLng, around int, limit int) ([]SearchResp, error) {
	query := fmt.Sprintf(`[out:json];
	(
	  node["name"]["tourism"="hotel"](around:%v,%v, %v);
	);
	out body;`, around, p.Lat, p.Lng)
	c := overpass.New()
	resp, err := c.Query(query)
	if err != nil {
		return []SearchResp{}, err
	}

	res := []SearchResp{}
	var count int
	for _, n := range resp.Nodes {
		if count == limit {
			break
		}
		res = append(res, SearchResp{Name: n.Tags["name"], P: s2.LatLngFromDegrees(n.Lat, n.Lon)})
		count++
	}

	return res, nil

}

func GetCafes(p s2.LatLng, around int, limit int) ([]SearchResp, error) {

	ar := fmt.Sprintf(`(around:%v,%v, %v)`, around, p.Lat, p.Lng)
	query := fmt.Sprintf(`
	[out:json];
(
  node["name"]["amenity"="restaurant"]%v;
  node["name"]["amenity"="cafe"]%v;
);
out body;
	`, ar, ar)

	c := overpass.New()
	resp, err := c.Query(query)
	if err != nil {
		return []SearchResp{}, err
	}

	res := []SearchResp{}
	var count int
	for _, n := range resp.Nodes {
		if count == limit {
			break
		}
		res = append(res, SearchResp{Name: n.Tags["name"], P: s2.LatLngFromDegrees(n.Lat, n.Lon)})
		count++
	}

	return res, nil

}
