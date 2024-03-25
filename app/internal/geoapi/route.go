package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/geo/s2"
)

var routeBaseUrl = "http://router.project-osrm.org/route/v1/driving/"

type routeResp struct {
	Routes []struct {
		Legs []struct {
			Steps []struct {
				Intersections []struct {
					Location []float64 `json:"location"`
				} `json:"intersections"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
}

func (r routeResp) Points() []s2.LatLng {
	res := []s2.LatLng{}

	for _, route := range r.Routes {
		for _, leg := range route.Legs {
			for _, step := range leg.Steps {
				for _, intersect := range step.Intersections {
					lat := intersect.Location[0]
					lng := intersect.Location[1]
					res = append(res, s2.LatLngFromDegrees(lat, lng))
				}
			}
		}
	}

	return res

}

func getRouteUrl(coords []s2.LatLng) string {
	steps := make([]string, len(coords))
	for i, c := range coords {
		steps[i] = fmt.Sprintf(`%v,%v`, c.Lat, c.Lng)
	}
	return fmt.Sprintf(`%v%v?overview=false&steps=true`, routeBaseUrl, strings.Join(steps, ";"))
}

func GetRoute(points []s2.LatLng) ([]s2.LatLng, error) {
	resp, err := http.Get(getRouteUrl(points))
	if err != nil {
		return []s2.LatLng{}, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []s2.LatLng{}, err
	}

	// fmt.Println(string(data))

	var r routeResp
	err = json.Unmarshal(data, &r)
	if err != nil {
		return []s2.LatLng{}, err
	}

	return r.Points(), nil
}
