package domain

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Accommodation struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner       User               `bson:"owner,omitempty" json:"owner"`
	Name        string             `bson:"name" json:"name"`
	Location    string             `bson:"location" json:"location"`
	MinGuestNum int                `bson:"minGuestNum" json:"minGuestNum"`
	MaxGuestNum int                `bson:"maxGuestNum" json:"maxGuestNum"`
	Amenities   []string           `bson:"amenities,omitempty" json:"amenities"`
}

type Accommodations []*Accommodation

func (a Accommodations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

func (u User) Equals(user User) bool {
	return u.Id.String() == user.Id.String()
}

func (a Accommodation) Of(user User) bool {
	return a.Owner.Equals(user)
}

//type AccommodationRepository interface {
//	Get(id primitive.ObjectID) (Accommodation, error)
//	GetAll() ([]Accommodation, error)
//	GetByUser(user User) ([]Accommodation, error)
//	Create(Accommodation Accommodation) (Accommodation, error)
//	Update(Accommodation Accommodation) error
//}
