package client

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserData struct {
	ID        primitive.ObjectID
	Name      *string
	Last_Name *string
	Email     string
	Username  string
	Address   *string
	Role      string
}
