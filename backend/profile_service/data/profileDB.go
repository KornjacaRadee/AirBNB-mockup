package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func CreateProfile(client *mongo.Client, profile *Profile) error {
	profileCollection := client.Database("profileDB").Collection("profiles")

	//currentTime := time.Now()
	//profile.Created_On = currentTime.Format(time.RFC3339)
	//profile.Updated_On = currentTime.Format(time.RFC3339)

	result, err := profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}

	profile.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func GetAllProfiles(client *mongo.Client) (Profiles, error) {
	profileCollection := client.Database("profileDB").Collection("profiles")

	cursor, err := profileCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var profiles Profiles
	for cursor.Next(context.TODO()) {
		var profile Profile
		if err := cursor.Decode(&profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func GetProfileByID(client *mongo.Client, profileID primitive.ObjectID) (*Profile, error) {
	profileCollection := client.Database("profileDB").Collection("profiles")

	var profile Profile
	err := profileCollection.FindOne(context.TODO(), bson.M{"_id": profileID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func UpdateProfile(client *mongo.Client, profileID primitive.ObjectID, updatedProfile *Profile) error {
	profileCollection := client.Database("profileDB").Collection("profiles")

	//updatedProfile.Updated_On = time.Now().Format(time.RFC3339)

	result, err := profileCollection.ReplaceOne(context.TODO(), bson.M{"_id": profileID}, updatedProfile)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("profile not found")
	}

	return nil
}

func DeleteProfile(client *mongo.Client, profileID primitive.ObjectID) error {
	profileCollection := client.Database("profileDB").Collection("profiles")

	result, err := profileCollection.DeleteOne(context.TODO(), bson.M{"_id": profileID})
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("profile not found")
	}

	return nil
}

func GetProfileByEmail(client *mongo.Client, email string) (*Profile, error) {
	profileCollection := client.Database("authDB").Collection("profiles")

	// Create a filter for the email
	filter := bson.D{{"email", email}}

	// Find the user in the database
	var profile Profile
	err := profileCollection.FindOne(context.TODO(), filter).Decode(&profile)
	if err != nil {
		// Handle the error (e.g., user not found)
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return &profile, nil
}
