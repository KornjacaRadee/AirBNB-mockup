package client

import "go.mongodb.org/mongo-driver/bson/primitive"

type AccommodationData struct {
	Id          primitive.ObjectID `json:"id"`
	Owner       User               `json:"owner"`
	Name        string             `json:"name"`
	Location    string             `json:"location"`
	MinGuestNum int                `json:"minGuestNum"`
	MaxGuestNum int                `json:"maxGuestNum"`
	Amenities   []string           `json:"amenities"`
}

type User struct {
	Id primitive.ObjectID `json:"id"`
}
