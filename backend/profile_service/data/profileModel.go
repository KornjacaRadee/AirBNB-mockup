// profile.go

package data

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Profile struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	First_Name *string            `bson:"first_name" json:"name" validate:"required"`
	Last_Name  *string            `bson:"last_name" json:"last_name"`
	Username   *string            `bson:"username" json:"username" validate:"required"`
	Email      string             `bson:"email" json:"email" validate:"required,email"`
	Address    *string            `bson:"address" json:"address"`
	Created_On string             `bson:"created_on" json:"created_On"`
	Updated_On string             `bson:"updated_on" json:"updated_On"`
	Deleted_On string             `bson:"deleted_on" json:"deleted_On"`
}

type Profiles []*Profile

func (p *Profiles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Profile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Profile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
