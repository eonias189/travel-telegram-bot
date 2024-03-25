package geoapi

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/utils"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/r1"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

func ConvertToBytes(img image.Image) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	err := png.Encode(buf, img)
	if err != nil {
		return []byte{}, err
	}

	return io.ReadAll(buf)
}

func GetRouteImg(points []s2.LatLng) (image.Image, error) {
	ctx := sm.NewContext()
	ctx.SetSize(800, 600)

	route, err := GetRoute(points)
	if err != nil {
		return nil, err
	}

	lats := make([]float64, len(points))
	longs := make([]float64, len(points))

	for i, p := range points {
		lats[i] = float64(p.Lat)
		longs[i] = float64(p.Lng)
		ctx.AddObject(sm.NewMarker(p, color.Opaque, 16))
	}

	rect := s2.Rect{Lat: r1.Interval{Lo: utils.Min(lats) - 0.01, Hi: utils.Max(lats) + 0.01}, Lng: s1.Interval{Lo: utils.Min(longs) - 0.01, Hi: utils.Max(longs) + 0.01}}
	ctx.AddObject(sm.NewPath(route, color.RGBA{R: 45, G: 48, B: 229, A: 100}, 10))
	ctx.SetBoundingBox(rect)
	return ctx.Render()
}
