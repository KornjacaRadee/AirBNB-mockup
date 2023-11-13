package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func registerUser(client *mongo.Client, user *User) error {
	userCollection := client.Database("your_database").Collection("users")
	_, err := userCollection.InsertOne(context.TODO(), user)
	return err
}

func loginUser(client *mongo.Client, email, password string) (*User, error) {
	userCollection := client.Database("your_database").Collection("users")

	var user User
	err := userCollection.FindOne(context.TODO(), bson.D{{"email", email}, {"password", password}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getAllUsers(client *mongo.Client) (Users, error) {
	userCollection := client.Database("your_database").Collection("users")

	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users Users
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func deleteUser(client *mongo.Client, userID primitive.ObjectID) error {
	userCollection := client.Database("your_database").Collection("users")

	_, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", userID}})
	return err
}
