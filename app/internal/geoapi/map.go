package geoapi

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"

	sm "github.com/flopp/go-staticmaps"
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

	for _, p := range points {
		ctx.AddObject(sm.NewMarker(p, color.Opaque, 16))
	}
	ctx.AddObject(sm.NewPath(route, color.RGBA{R: 45, G: 48, B: 229, A: 100}, 10))
	return ctx.Render()
}
