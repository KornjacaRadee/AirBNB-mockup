package data

import (
	"auth-service/config"
	"bufio"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

var logger = config.NewLogger("./logging/log.log")

func RegisterUser(client *mongo.Client, user *User) error {
	userCollection := client.Database("authDB").Collection("users")

	// Create unique index on email field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		logger.Errorf("Error creating unique index: %v", err)
		return err
	}

	// Try to insert the user
	_, err = userCollection.InsertOne(context.TODO(), user)

	// Check for duplicate key error
	if writeException, ok := err.(mongo.WriteException); ok {
		for _, writeError := range writeException.WriteErrors {
			if writeError.Code == 11000 { // Duplicate key error code
				logger.Warnf("Email '%s' is already registered", user.Email)
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
		logger.Errorf("Error finding user: %v", err)
		return nil, err
	}

	logger.Debugf("Retrieved hashed password for user '%s'", user.Email)

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Warnf("Password mismatch for user '%s'", user.Email)
		return nil, err // Passwords do not match
	}

	logger.Infof("User '%s' successfully logged in", user.Email)
	return &user, nil
}
func UpdatePassword(client *mongo.Client, userID primitive.ObjectID, newPassword string) error {
	// Hash the new password
	logger.Debugf("Updating password for user with ID '%s'", userID.Hex())
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		logger.Errorf("Error hashing password: %v", err)
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
	if err != nil {
		logger.Errorf("Error updating password: %v", err)
		return err
	}

	logger.Infof("Password updated successfully for user with ID '%s'", userID.Hex())
	return nil
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

func GetUserByEmail(client *mongo.Client, email string) (*User, error) {
	userCollection := client.Database("authDB").Collection("users")

	// Create a filter for the email
	filter := bson.D{{"email", email}}

	// Find the user in the database
	var user User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		// Handle the error (e.g., user not found)
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return &user, nil
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
		logger.Errorf("Error opening blacklist file: %v", err)
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == password {
			logger.Warnf("Password found in blacklist: '%s'", password)
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("Error scanning blacklist: %v", err)
		return false, err
	}

	logger.Debugf("Password '%s' not found in blacklist", password)
	return true, nil
}
