package app

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

func handleLocations(opts AppHandlerOptions, locationService LocationService) {

	var getInt64 = func(query url.Values, key string) (int64, error) {
		str := query.Get(key)
		if str == "" {
			return 0, ErrInternal
		}

		rInt, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}

		return int64(rInt), nil
	}

	var renderLocations = func(tripId int64, ctx *tgapi.Context) error {
		locations, err := locationService.GetAll(tripId)
		if err != nil {
			return err
		}
		return ctx.SendMessage(msgtempl.LocationsMessage(ctx.SenderID(), tripId, locations))
	}

	var renderLocation = func(tripId int64, order int, ctx *tgapi.Context) error {
		location, err := locationService.Get(tripId, order)
		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.LocationMessage(tripId, ctx.SenderID(), location))
	}

	opts.CallbackRouter.Handle("locations", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		return renderLocations(tripId, ctx)
	})

	opts.CallbackRouter.Handle("new-location", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`location-name-input?tripId=%v`, tripId))
		return ctx.SendString("введи локацию")
	})

	opts.ContextRouter.Handle("location-name-input", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(opts.Dcqp.GetDialogContextQuery(ctx), "tripId")
		if err != nil {
			return err
		}

		location := ctx.Update.Message.Text
		if !geoapi.CheckLocation(location) {
			return ctx.SendString(fmt.Sprintf(`%v: не найдено`, location))
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`location-start-input?tripId=%v&name=%v`, tripId, location))
		return ctx.SendString("введи дату прибытия в формате день.месяц.год/часы:минуты")
	})

	opts.ContextRouter.Handle("location-start-input", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(opts.Dcqp.GetDialogContextQuery(ctx), "tripId")
		if err != nil {
			return err
		}

		name := opts.Dcqp.GetDialogContextQuery(ctx).Get("name")
		if name == "" {
			return ErrInternal
		}

		text := ctx.Update.Message.Text
		layout := `02.01.2006/15:04`
		t, err := time.Parse(layout, text)
		if err != nil {
			return ctx.SendString("неверный формат")
		}

		if time.Now().After(t) {
			return ctx.SendString("дата прибытия должна быть в будущем времени")
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`location-end-input?tripId=%v&name=%v&start=%v`, tripId, name, t.UnixNano()))
		return ctx.SendString("введи дату отбытия в формате день.месяц.год/часы:минуты")
	})

	opts.ContextRouter.Handle("location-end-input", func(ctx *tgapi.Context) error {
		query := opts.Dcqp.GetDialogContextQuery(ctx)
		tripId, err := getInt64(query, "tripId")
		if err != nil {
			return err
		}

		name := query.Get("name")
		if name == "" {
			return ErrInternal
		}

		start, err := getInt64(query, "start")
		if err != nil {
			return err
		}

		end := ctx.Update.Message.Text
		layout := `02.01.2006/15:04`
		t, err := time.Parse(layout, end)
		if err != nil {
			return ctx.SendString("неверный формат")
		}
		if t.UnixNano() < start {
			return ctx.SendString("дата отбытия должна быть позже даты прибытия")
		}

		locations, err := locationService.GetAll(tripId)
		if err != nil {
			return err
		}

		location := service.Location{
			Order:     len(locations) + 1,
			Name:      name,
			StartTime: start,
			EndTime:   t.UnixNano(),
		}

		err = locationService.Add(tripId, location)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")

		return renderLocations(tripId, ctx)

	})

	opts.CallbackRouter.Handle("location", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		orderInt64, err := getInt64(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		order := int(orderInt64)
		return renderLocation(tripId, order, ctx)
	})

	opts.CallbackRouter.Handle("delete-location", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		orderInt64, err := getInt64(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		order := int(orderInt64)
		err = locationService.Delete(tripId, order)
		if err != nil {
			return err
		}

		return renderLocations(tripId, ctx)
	})

	opts.CallbackRouter.Handle("change-location-start", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		orderInt64, err := getInt64(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`change-location-start?tripId=%v&order=%v`, tripId, int(orderInt64)))
		return ctx.SendString("введи новую дату прибытия в формате день.месяц.год/часы:минуты")
	})

	opts.ContextRouter.Handle("change-location-start", func(ctx *tgapi.Context) error {
		query := opts.Dcqp.GetDialogContextQuery(ctx)
		tripId, err := getInt64(query, "tripId")
		if err != nil {
			return err
		}

		order, err := getInt64(query, "order")
		if err != nil {
			return err
		}

		loc, err := locationService.Get(tripId, int(order))
		if err != nil {
			return err
		}

		start := ctx.Update.Message.Text
		layout := `02.01.2006/15:04`
		t, err := time.Parse(layout, start)
		if err != nil {
			return ctx.SendString("неверный формат")
		}
		if time.Now().After(t) {
			return ctx.SendString("дата прибытия должна быть в будущем времени")
		}
		loc.StartTime = t.UnixNano()

		err = locationService.Set(tripId, int(order), loc)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return renderLocation(tripId, int(order), ctx)
	})

	opts.CallbackRouter.Handle("change-location-end", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		orderInt64, err := getInt64(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`change-location-end?tripId=%v&order=%v`, tripId, int(orderInt64)))
		return ctx.SendString("введи новую дату отбытия в формате день.месяц.год/часы:минуты")
	})

	opts.ContextRouter.Handle("change-location-end", func(ctx *tgapi.Context) error {
		query := opts.Dcqp.GetDialogContextQuery(ctx)
		tripId, err := getInt64(query, "tripId")
		if err != nil {
			return err
		}

		order, err := getInt64(query, "order")
		if err != nil {
			return err
		}

		loc, err := locationService.Get(tripId, int(order))
		if err != nil {
			return err
		}

		end := ctx.Update.Message.Text
		layout := `02.01.2006/15:04`
		t, err := time.Parse(layout, end)
		if err != nil {
			return ctx.SendString("неверный формат")
		}
		if t.UnixNano() < loc.StartTime {
			return ctx.SendString("дата отбытия должна быть позже даты прибытия")
		}

		loc.EndTime = t.UnixNano()

		err = locationService.Set(tripId, int(order), loc)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return renderLocation(tripId, int(order), ctx)
	})

	opts.CallbackRouter.Handle("change-location-name", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		orderInt64, err := getInt64(ctx.CallbackQuery(), "order")
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, fmt.Sprintf(`change-location-name?tripId=%v&order=%v`, tripId, int(orderInt64)))
		return ctx.SendString("введи новую локацию")
	})

	opts.ContextRouter.Handle("change-location-name", func(ctx *tgapi.Context) error {
		query := opts.Dcqp.GetDialogContextQuery(ctx)
		tripId, err := getInt64(query, "tripId")
		if err != nil {
			return err
		}

		order, err := getInt64(query, "order")
		if err != nil {
			return err
		}

		loc, err := locationService.Get(tripId, int(order))
		if err != nil {
			return err
		}

		location := ctx.Update.Message.Text
		if !geoapi.CheckLocation(location) {
			return ctx.SendString(fmt.Sprintf(`%v: не найдено`, location))
		}

		loc.Name = location

		err = locationService.Set(tripId, int(order), loc)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return renderLocation(tripId, int(order), ctx)
	})
}
