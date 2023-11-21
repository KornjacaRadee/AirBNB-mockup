package domain

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type AvailabilityPeriod struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Accommodation Accommodation      `bson:"accommodation,omitempty" json:"accommodation"`
	StartDate     string             `bson:"startDate" json:"startDate"`
	EndDate       string             `bson:"endDate" json:"endDate"`
	Price         int                `bson:"price" json:"price"`
}

type AvailabilityPeriods []*AvailabilityPeriod

func (a AvailabilityPeriods) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AvailabilityPeriod) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

type Accommodation struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

//type AccommodationRepository interface {
//	Get(id primitive.ObjectID) (Accommodation, error)
//	GetAll() ([]Accommodation, error)
//	GetByUser(user User) ([]Accommodation, error)
//	Create(Accommodation Accommodation) (Accommodation, error)
//	Update(Accommodation Accommodation) error
//}
