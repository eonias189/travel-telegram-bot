package app

import (
	"fmt"
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/utils"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
	"github.com/golang/geo/s2"
)

func handleGeoTools(opts AppHandlerOptions, locationService LocationService) {
	opts.CallbackRouter.Handle("get-attractions", func(ctx *tgapi.Context) error {
		tripId, err := utils.GetInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		order, err := utils.GetInt(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		loc, err := locationService.Get(tripId, order)
		if err != nil {
			return err
		}

		p := s2.LatLngFromDegrees(loc.Lat, loc.Lng)
		address, err := geoapi.GetAddress(loc.Name)
		if err != nil {
			return err
		}

		attractions, err := geoapi.GetAttractions(p, 10000)
		if err != nil {
			return err
		}

		attractions = utils.Map(attractions, func(attr string) string {
			return fmt.Sprintf(`%v (%v, %v)`, attr, address.Country, address.City)
		})

		return ctx.SendString(fmt.Sprintf("Достопримечательности рядом:\n%v", strings.Join(attractions, ",\n")))

	})
}
