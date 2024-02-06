package client

import (
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReservationData struct {
	AvailabilityPeriodId gocql.UUID
	Id                   gocql.UUID
	StartDate            time.Time
	EndDate              time.Time
	AccommodationId      primitive.ObjectID
	HostId               primitive.ObjectID
	GuestId              primitive.ObjectID
	GuestNum             int
	Price                int
}

type ReservationsData []*ReservationData

type SearchReq struct {
	StartDate       time.Time
	EndDate         time.Time
	AccommodationId primitive.ObjectID
}

type SearchReqs []*SearchReq
