package domain

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type HostRating struct {
	Id      gocql.UUID
	GuestId primitive.ObjectID
	HostId  primitive.ObjectID
	Time    time.Time
	Rating  int
}

type AccommodationRating struct {
	Id              gocql.UUID
	GuestId         primitive.ObjectID
	HostId          primitive.ObjectID
	AccommodationId primitive.ObjectID
	Time            time.Time
	Rating          int
}

type HostRatings []*HostRating
type AccommodationRatings []*AccommodationRating

func (a HostRatings) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *HostRatings) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a AccommodationRatings) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AccommodationRatings) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a HostRating) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *HostRating) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a AccommodationRating) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AccommodationRating) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
