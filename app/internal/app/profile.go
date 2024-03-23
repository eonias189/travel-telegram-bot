package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/router"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type UserService interface {
	Get(id int64) (service.User, error)
	Set(id int64, user service.User) error
}

type DialogContextSetter interface {
	SetDialogContext(ctx *tgapi.Context, dialogContext string)
}

func handleProfile(contextrouter *router.Router, callbackRouter *router.Router, userService UserService, dcs DialogContextSetter) {
	callbackRouter.Handle("profile", func(ctx *tgapi.Context) error {

		sender := ctx.Update.SentFrom()
		user, err := userService.Get(sender.ID)
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			return err
		}

		msgtext := msgtempl.ProfileMessage(sender.UserName, user)
		return ctx.SendWithInlineKeyboard(msgtext, msgtempl.ProfileButtons())
	})

	callbackRouter.Handle("change-age", func(ctx *tgapi.Context) error {
		dcs.SetDialogContext(ctx, "change-age")
		return ctx.SendString("введи новый возраст")
	})

	contextrouter.Handle("change-age", func(ctx *tgapi.Context) error {
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

		dcs.SetDialogContext(ctx, "")
		return ctx.SendWithInlineKeyboard(msgtempl.ProfileMessage(sender.UserName, user), msgtempl.ProfileButtons())
	})

	callbackRouter.Handle("change-location", func(ctx *tgapi.Context) error {
		dcs.SetDialogContext(ctx, "change-location")
		return ctx.SendString("введи новое местоположение")
	})

	contextrouter.Handle("change-location", func(ctx *tgapi.Context) error {
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

		dcs.SetDialogContext(ctx, "")
		return ctx.SendWithInlineKeyboard(msgtempl.ProfileMessage(sender.UserName, user), msgtempl.ProfileButtons())
	})

	callbackRouter.Handle("change-bio", func(ctx *tgapi.Context) error {
		dcs.SetDialogContext(ctx, "change-bio")
		return ctx.SendString("введи новое bio")
	})

	contextrouter.Handle("change-bio", func(ctx *tgapi.Context) error {
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

		dcs.SetDialogContext(ctx, "")
		return ctx.SendWithInlineKeyboard(msgtempl.ProfileMessage(sender.UserName, user), msgtempl.ProfileButtons())
	})
}
