package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//func RegisterUser(client *mongo.Client, user *User) error {
//	userCollection := client.Database("mongodb").Collection("users")
//
//	// Create unique index on email field
//	indexModel := mongo.IndexModel{
//		Keys:    bson.D{{"email", 1}},
//		Options: options.Index().SetUnique(true),
//	}
//	_, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
//	if err != nil {
//		return err
//	}
//
//	// Try to insert the user
//	_, err = userCollection.InsertOne(context.TODO(), user)
//
//	// Check for duplicate key error
//	if writeException, ok := err.(mongo.WriteException); ok {
//		for _, writeError := range writeException.WriteErrors {
//			if writeError.Code == 11000 { // Duplicate key error code
//				return fmt.Errorf("email '%s' is already registered", user.Email)
//			}
//		}
//	}
//
//	return err
//}

func RegisterUser(client *mongo.Client, user *User) error {
	userCollection := client.Database("mongodb").Collection("users")

	// Create unique index on email field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	// Try to insert the user
	_, err = userCollection.InsertOne(context.TODO(), user)

	// Check for duplicate key error
	if writeException, ok := err.(mongo.WriteException); ok {
		for _, writeError := range writeException.WriteErrors {
			if writeError.Code == 11000 { // Duplicate key error code
				return fmt.Errorf("email '%s' is already registered", user.Email)
			}
		}
	}

	return err
}

func LoginUser(client *mongo.Client, email, password string) (*User, error) {
	userCollection := client.Database("mongodb").Collection("users")

	var user User
	err := userCollection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	log.Printf("Retrieved hashed password: %s", user.Password)
	log.Printf("Entered password: %s", password)
	// Verify the password
	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing passwords: %v", err)
		return nil, err // Passwords do not match
	}

	// Passwords match, return the user
	return &user, nil
}

func GetAllUsers(client *mongo.Client) (Users, error) {
	userCollection := client.Database("mongodb").Collection("users")

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

func DeleteUser(client *mongo.Client, userID primitive.ObjectID) error {
	userCollection := client.Database("mongodb").Collection("users")

	_, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", userID}})
	return err
}
