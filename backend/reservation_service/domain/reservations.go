package domain

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type AvailabilityPeriodByAccommodation struct {
	AccommodationId primitive.ObjectID
	HostId          primitive.ObjectID
	Id              gocql.UUID
	StartDate       time.Time
	EndDate         time.Time
	Price           int
	IsPricePerGuest bool
}

type ReservationByAvailabilityPeriod struct {
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

type AvailabilityPeriodsByAccommodation []*AvailabilityPeriodByAccommodation
type ReservationsByAvailabilityPeriod []*ReservationByAvailabilityPeriod

func (a AvailabilityPeriodsByAccommodation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AvailabilityPeriodsByAccommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a ReservationsByAvailabilityPeriod) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *ReservationsByAvailabilityPeriod) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a AvailabilityPeriodByAccommodation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AvailabilityPeriodByAccommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a ReservationByAvailabilityPeriod) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *ReservationByAvailabilityPeriod) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
