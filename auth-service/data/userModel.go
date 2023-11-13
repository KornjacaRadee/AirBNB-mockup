package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	First_Name *string            `bson:"first_name" json:"name" validate:"required"`
	Last_Name  *string            `bson:"last_name" json:"description"`
	Email      *string            `bson:"email" json:"email" validate:"gt=0"`
	Address    *string            `bson:"address" json:"address"`
	Created_On string             `bson:"created_on" json:"created_On"`
	Updated_On string             `bson:"updated_on" json:"updated_On"`
	Deleted_On string             `bson:"deleted_on" json:"deleted_On"`
}

type Users []*User

// Functions to encode and decode products to json and from json.
// If we inject an interfacchocoe we achieve dependancy injection, meaning that anything that implements this interface can be passed down
// For us it will be ResponseWriter, but it also may be a file writer or something similar.
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (p *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
