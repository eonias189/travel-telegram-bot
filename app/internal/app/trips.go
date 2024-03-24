package app

import (
	"errors"
	"fmt"
	"strconv"

	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

func handleTrips(opts AppHandlerOptions, userService UserService, tripService TripService) {

	var renderTrips = func(ctx *tgapi.Context) error {
		user, err := userService.Get(ctx.SenderID())

		if errors.Is(err, service.ErrNotFound) {
			return ctx.SendString(`укажите информацию о себе в профиле для использования функции "путешествия"`)
		}

		if err != nil {
			return err
		}

		trips, err := tripService.GetAll(user.Trips)
		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.TripsMessage(ctx.SenderID(), trips))
	}

	opts.CallbackRouter.Handle("trips", func(ctx *tgapi.Context) error {
		return renderTrips(ctx)
	})

	opts.CallbackRouter.Handle("new-trip", func(ctx *tgapi.Context) error {
		opts.Dcs.SetDialogContext(ctx, "trip-name-input")
		return ctx.SendString("введи название путешествия")
	})

	opts.ContextRouter.Handle("trip-name-input", func(ctx *tgapi.Context) error {
		name := ctx.Update.Message.Text
		if tripService.ExistsName(name) {
			return ctx.SendString("введённое имя путешествия уже занято")
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf("trip-description-input?name=%v", name))
		return ctx.SendString("введи описание путешествия")
	})

	opts.ContextRouter.Handle("trip-description-input", func(ctx *tgapi.Context) error {
		desc := ctx.Update.Message.Text
		query := opts.Dcqp.GetDialogContextQuery(ctx)
		name := query.Get("name")

		if name == "" {
			return ErrInternal
		}

		id := NewId()
		for tripService.Exists(id) {
			id = NewId()
		}

		trip := service.Trip{
			Id:          id,
			Name:        name,
			Description: desc,
			Creator:     ctx.SenderID(),
			Members:     []int64{ctx.SenderID()},
		}

		err := tripService.Set(id, trip)
		if err != nil {
			return err
		}

		err = userService.AddTrip(ctx.SenderID(), trip.Id)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return renderTrips(ctx)

	})

	opts.CallbackRouter.Handle("delete-trip", func(ctx *tgapi.Context) error {
		query := ctx.CallbackQuery()
		idStr := query.Get("id")

		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		id := int64(idInt)

		trip, err := tripService.Get(id)
		if err != nil {
			return err
		}

		err = tripService.Delete(id)
		if err != nil {

			return err
		}

		for _, memberId := range trip.Members {
			err = userService.DeleteTrip(memberId, id)
			if err != nil {
				opts.Logger.Error(err.Error())
			}
		}

		return renderTrips(ctx)

	})

	opts.CallbackRouter.Handle("trip", func(ctx *tgapi.Context) error {
		query := ctx.CallbackQuery()
		idStr := query.Get("id")
		if idStr == "" {
			return ErrInternal
		}

		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		id := int64(idInt)
		trip, err := tripService.Get(id)
		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.TripMessage(ctx.SenderID(), trip))

	})
}
