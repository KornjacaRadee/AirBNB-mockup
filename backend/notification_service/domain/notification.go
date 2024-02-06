package domain

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type Notification struct {
	Id   string    `json:"id"`
	Host User      `json:"host"`
	Text string    `son:"text"`
	Time time.Time `json:"time"`
}

type Notifications []*Notification

func (n Notifications) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(n)
}

func (n *Notification) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(n)
}

func (n *Notification) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(n)
}

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

func (u User) Equals(user User) bool {
	return u.Id.String() == user.Id.String()
}

func (n Notification) Of(user User) bool {
	return n.Host.Equals(user)
}
