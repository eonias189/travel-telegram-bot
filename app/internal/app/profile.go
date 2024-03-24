package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

func handleProfile(opts AppHandlerOptions, userService UserService) {
	opts.CallbackRouter.Handle("profile", func(ctx *tgapi.Context) error {

		sender := ctx.Update.SentFrom()
		user, err := userService.Get(sender.ID)
		if errors.Is(err, service.ErrNotFound) {
			return ctx.SendMessage(msgtempl.ProfileMsg(sender.ID, sender.UserName, service.User{
				Age:      -1,
				Location: "не указано",
				Bio:      "не указано",
			}))
		}

		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.ProfileMsg(sender.ID, sender.UserName, user))
	})

	opts.CallbackRouter.Handle("change-age", func(ctx *tgapi.Context) error {
		opts.Dcs.SetDialogContext(ctx, "change-age")
		return ctx.SendString("введи новый возраст")
	})

	opts.ContextRouter.Handle("change-age", func(ctx *tgapi.Context) error {
		newAge := ctx.Update.Message.Text
		newAgetInt, err := strconv.Atoi(newAge)
		if err != nil {
			return ctx.SendString("возраст некорректен")
		}

		sender := ctx.Update.SentFrom()
		user, err := userService.Get(sender.ID)
		if err != nil && !errors.Is(service.ErrNotFound, err) {
			return err
		}

		user.Age = newAgetInt
		err = userService.Set(sender.ID, user)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return ctx.SendMessage(msgtempl.ProfileMsg(sender.ID, sender.UserName, user))
	})

	opts.CallbackRouter.Handle("change-location", func(ctx *tgapi.Context) error {
		opts.Dcs.SetDialogContext(ctx, "change-location")
		return ctx.SendString("введи новое местоположение")
	})

	opts.ContextRouter.Handle("change-location", func(ctx *tgapi.Context) error {
		sender := ctx.Update.SentFrom()

		newLocation := ctx.Update.Message.Text
		user, err := userService.Get(sender.ID)
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			return err
		}

		if !geoapi.CheckLocation(newLocation) {
			return ctx.SendString(fmt.Sprintf("%v: не найдено", newLocation))
		}

		user.Location = newLocation
		err = userService.Set(sender.ID, user)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return ctx.SendMessage(msgtempl.ProfileMsg(sender.ID, sender.UserName, user))
	})

	opts.CallbackRouter.Handle("change-bio", func(ctx *tgapi.Context) error {
		opts.Dcs.SetDialogContext(ctx, "change-bio")
		return ctx.SendString(`введи новое "о себе"`)
	})

	opts.ContextRouter.Handle("change-bio", func(ctx *tgapi.Context) error {
		sender := ctx.Update.SentFrom()

		newBio := ctx.Update.Message.Text
		user, err := userService.Get(sender.ID)
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			return err
		}

		user.Bio = newBio
		err = userService.Set(sender.ID, user)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")
		return ctx.SendMessage(msgtempl.ProfileMsg(sender.ID, sender.UserName, user))
	})
}
