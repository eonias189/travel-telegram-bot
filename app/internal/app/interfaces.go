package app

import (
	"log/slog"
	"net/url"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/router"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type AppHandlerOptions struct {
	CommandRouter  *router.Router
	ContextRouter  *router.Router
	CallbackRouter *router.Router
	Dcs            DialogContextSetter
	Dcqp           DialogContextQueryProvider
	Logger         *slog.Logger
}

type DialogContextSetter interface {
	SetDialogContext(ctx *tgapi.Context, dialogContext string)
}

type DialogContextQueryProvider interface {
	GetDialogContextQuery(ctx *tgapi.Context) url.Values
}

type UserService interface {
	Get(id int64) (service.User, error)
	Set(id int64, user service.User) error
	AddTrip(userId, tripId int64) error
	DeleteTrip(userId, tripId int64) error
}

type TripService interface {
	Get(id int64) (service.Trip, error)
	Set(id int64, trip service.Trip) error
	Exists(id int64) bool
	ExistsName(name string) bool
	GetAll(ids []int64) ([]service.Trip, error)
	Delete(id int64) error
}

type LocationService interface {
	Get(tripId int64, order int) (service.Location, error)
	Set(tripId int64, order int, location service.Location) error
	GetAll(tripId int64) ([]service.Location, error)
	Add(tripId int64, location service.Location) error
	Delete(tripId int64, order int) error
}
