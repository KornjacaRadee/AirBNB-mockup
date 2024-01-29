package client

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotificationData struct {
	Host User      `bson:"host,omitempty" json:"host"`
	Text string    `bson:"text" json:"text"`
	Time time.Time `bson:"time" json:"time"`
}

type User struct {
	Id primitive.ObjectID `json:"id"`
}
