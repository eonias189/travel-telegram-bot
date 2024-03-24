package app

import (
	"fmt"
	"strconv"

	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

func handleTrips(opts AppHandlerOptions, userService UserService, tripService TripService) {

	var renderTrips = func(ctx *tgapi.Context) error {
		user, err := userService.Get(ctx.SenderID())
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
		opts.Dcs.SetDialogContext(ctx, "name-input")
		return ctx.SendString("введи название путешествия")
	})

	opts.ContextRouter.Handle("name-input", func(ctx *tgapi.Context) error {
		name := ctx.Update.Message.Text
		if tripService.ExistsName(name) {
			return ctx.SendString("введённое имя путешествия уже занято")
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf("description-input?name=%v", name))
		return ctx.SendString("введи описание путешествия")
	})

	opts.ContextRouter.Handle("description-input", func(ctx *tgapi.Context) error {
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

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		err = tripService.Delete(int64(id))
		if err != nil {

			return err
		}

		err = userService.DeleteTrip(ctx.SenderID(), int64(id))
		if err != nil {
			return err
		}

		return renderTrips(ctx)

	})
}
