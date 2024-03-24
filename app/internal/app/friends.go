package app

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"

	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/utils"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SharePayload struct {
	TripId string `json:"tripId"`
}

var secretKey = "very very secret"

func handleFriends(opts AppHandlerOptions, userService UserService, tripService TripService) {

	var getInt64 = func(query url.Values, key string) (int64, error) {
		resStr := query.Get(key)
		if resStr == "" {
			return 0, ErrInternal
		}

		resInt, err := strconv.Atoi(resStr)
		if err != nil {
			return 0, ErrInternal
		}

		return int64(resInt), nil
	}

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

	var renderTrip = func(ctx *tgapi.Context, tripId int64) error {
		trip, err := tripService.Get(tripId)
		if err != nil {
			return err
		}

		return ctx.SendMessage(msgtempl.TripMessage(ctx.SenderID(), trip))
	}

	opts.CallbackRouter.Handle("share-trip", func(ctx *tgapi.Context) error {
		tripId, err := getInt64(ctx.CallbackQuery(), "tripId")
		if err != nil {
			return err
		}

		token, err := utils.GenerateJWT(SharePayload{TripId: fmt.Sprint(tripId)}, secretKey)
		if err != nil {
			return err
		}

		err = ctx.SendString(fmt.Sprintf(`Код для присоединения: %v`, token))
		if err != nil {
			return err
		}

		return renderTrip(ctx, tripId)
	})

	opts.CallbackRouter.Handle("join-to-trip", func(ctx *tgapi.Context) error {
		opts.Dcs.SetDialogContext(ctx, "join-token-input")
		return ctx.SendString("введи код для присоединения")
	})

	opts.ContextRouter.Handle("join-token-input", func(ctx *tgapi.Context) error {
		token := ctx.Update.Message.Text

		var payload SharePayload
		err := utils.ReadJWT(&payload, token, secretKey)
		if err != nil {
			return ctx.SendString("некорректный код")
		}

		tripIdInt, err := strconv.Atoi(payload.TripId)
		if err != nil {
			return err
		}
		tripId := int64(tripIdInt)

		trip, err := tripService.Get(tripId)
		if err != nil {
			return ctx.SendString("путешествие не найдено")
		}

		if slices.Contains(trip.Members, ctx.SenderID()) {
			return ctx.SendString("вы уже состоите в этом путешествии")
		}

		err = tripService.AddMember(tripId, ctx.SenderID())
		if err != nil {
			return err
		}

		err = userService.AddTrip(ctx.SenderID(), tripId)
		if err != nil {
			return err
		}

		opts.Dcs.SetDialogContext(ctx, "")

		go func() {
			for _, member := range trip.Members {
				msg := tgbotapi.NewMessage(member, fmt.Sprintf(`%v присоединился к путешествию "%v"`, ctx.Update.SentFrom().UserName, trip.Name))
				if _, err := ctx.Bot.Send(msg); err != nil {
					opts.Logger.Error(err.Error())
				}
			}
		}()

		return renderTrips(ctx)
	})
}
