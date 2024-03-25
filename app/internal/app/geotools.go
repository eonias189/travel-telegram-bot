package app

import (
	"fmt"
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/utils"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
	"github.com/golang/geo/s2"
)

func handleGeoTools(opts AppHandlerOptions, locationService LocationService) {

	var renderLocation = func(ctx *tgapi.Context, tripId int64, order int) error {
		location, err := locationService.Get(tripId, order)
		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.LocationMessage(tripId, ctx.SenderID(), location))
	}

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

		attractions, err := geoapi.GetAttractions(p, 10000, 10)
		if err != nil {
			return err
		}

		if len(attractions) == 0 {
			return ctx.SendString("достопримечательностей рядом не найдено")
		}

		attractionStrings := utils.Map(attractions, func(attr geoapi.SearchResp) string {
			addr, err := geoapi.GetAddress(attr.P)
			if err != nil {
				opts.Logger.Error(err.Error())
				return fmt.Sprintf(`%v`, attr.Name)
			}
			return fmt.Sprintf(`%v (%v, %v)`, attr.Name, addr.Country, addr.City)
		})

		err = ctx.SendString(fmt.Sprintf("Достопримечательности рядом:\n%v", strings.Join(attractionStrings, ",\n")))
		if err != nil {
			return err
		}

		return renderLocation(ctx, tripId, order)
	})

	opts.CallbackRouter.Handle("get-hotels", func(ctx *tgapi.Context) error {
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
		hotels, err := geoapi.GetHotels(p, 10000, 10)
		if err != nil {
			return err
		}

		if len(hotels) == 0 {
			return ctx.SendString("отелей рядом не найдено")
		}

		hotelStrings := utils.Map(hotels, func(h geoapi.SearchResp) string {
			addr, err := geoapi.GetAddress(h.P)
			if err != nil {
				opts.Logger.Error(err.Error())
				return fmt.Sprintf(`%v (адрес не найден)`, h.Name)
			}
			return fmt.Sprintf(`%v: г. %v ул. %v %v`, h.Name, addr.City, addr.Road, addr.HouseNumber)
		})

		text := fmt.Sprintf("Отели рядом:\n%v", strings.Join(hotelStrings, ",\n"))
		err = ctx.SendString(text)
		if err != nil {
			return err
		}

		return renderLocation(ctx, tripId, order)
	})

	opts.CallbackRouter.Handle("get-cafes", func(ctx *tgapi.Context) error {
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
		cafes, err := geoapi.GetCafes(p, 10000, 10)
		if err != nil {
			return err
		}

		if len(cafes) == 0 {
			return ctx.SendString("ресторанов и кафе рядом не найдено")
		}

		cafeStrings := utils.Map(cafes, func(c geoapi.SearchResp) string {
			addr, err := geoapi.GetAddress(c.P)
			if err != nil {
				opts.Logger.Error(err.Error())
				return fmt.Sprintf(`%v (адрес не найден)`, c.Name)
			}
			return fmt.Sprintf(`%v: г. %v ул. %v %v`, c.Name, addr.City, addr.Road, addr.HouseNumber)
		})

		text := fmt.Sprintf("Отели рядом:\n%v", strings.Join(cafeStrings, ",\n"))
		err = ctx.SendString(text)
		if err != nil {
			return err
		}

		return renderLocation(ctx, tripId, order)
	})
}
