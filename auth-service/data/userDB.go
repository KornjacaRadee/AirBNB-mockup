package data

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

/*func RegisterUser(client *mongo.Client, user *User) error {
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
}*/

func RegisterUser(client *mongo.Client, user *User) error {
	userCollection := client.Database("authDB").Collection("users")

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

// data/user.go

// ...

func GetUserByID(client *mongo.Client, userID primitive.ObjectID) (*User, error) {
	userCollection := client.Database("authDB").Collection("users")

	var user User
	err := userCollection.FindOne(context.TODO(), bson.D{{"_id", userID}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(client *mongo.Client, email, password string) (*User, error) {
	userCollection := client.Database("authDB").Collection("users")

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

	return &user, nil
}
func UpdatePassword(client *mongo.Client, userID primitive.ObjectID, newPassword string) error {
	// Hash the new password
	log.Printf("New password is: %s ", newPassword)
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the user's password in the database
	userCollection := client.Database("authDB").Collection("users")
	filter := bson.D{{"_id", userID}}
	update := bson.D{
		{"$set", bson.D{
			{"password", hashedPassword},
		}},
	}

	_, err = userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func GetAllUsers(client *mongo.Client) (Users, error) {
	userCollection := client.Database("authDB").Collection("users")

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
	userCollection := client.Database("authDB").Collection("users")

	_, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", userID}})
	return err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordInBlacklist(password string) (bool, error) {
	file, err := os.Open("blacklist/blacklist.txt")
	if err != nil {
		log.Printf("error while opening blacklist file: %v", err)
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == password {
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error while scanning blackist: %v", err)
		return false, err
	}

	return true, nil
}
