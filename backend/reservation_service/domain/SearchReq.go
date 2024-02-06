package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SearchReq struct {
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	AccommodationId primitive.ObjectID `json:"accommodation_id"`
}

type SearchReqs []*SearchReq
