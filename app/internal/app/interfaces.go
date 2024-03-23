package app

import (
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type DialogContextSetter interface {
	SetDialogContext(ctx *tgapi.Context, dialogContext string)
}

type UserService interface {
	Get(id int64) (service.User, error)
	Set(id int64, user service.User) error
}

type TripService interface {
	Get(id int64) (service.Trip, error)
	Set(id int64, trip service.Trip) error
}
