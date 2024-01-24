package client

import (
	"time"
)

type NotificationData struct {
	Host User      `bson:"host,omitempty" json:"host"`
	Text string    `bson:"text" json:"text"`
	Time time.Time `bson:"time" json:"time"`
}
