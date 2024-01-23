package client

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserData struct {
	Id       primitive.ObjectID `json:"id"`
	Name     string             `json:"name"`
	LastName string             `json:"last_name"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Address  string             `json:"address"`
	Role     string             `json:"role"`
}
