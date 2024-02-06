package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SearchReq struct {
	StartDate       time.Time
	EndDate         time.Time
	AccommodationId primitive.ObjectID
}

type SearchReqs []*SearchReq
