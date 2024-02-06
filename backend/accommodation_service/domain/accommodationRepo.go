package domain

import (
	"accommodation_service/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// NoSQL: AccommodationRepo struct encapsulating Mongo api client
type AccommodationRepo struct {
	cli    *mongo.Client
	logger *config.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *config.Logger) (*AccommodationRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &AccommodationRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (ar *AccommodationRepo) Disconnect(ctx context.Context) error {
	err := ar.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (ar *AccommodationRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := ar.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		ar.logger.Println(err)
	}

	// Print available databases
	databases, err := ar.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
	}
	fmt.Println(databases)
}

func (ar *AccommodationRepo) GetAll() (Accommodations, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()

	var accommodations Accommodations
	accommodationsCursor, err := accommodationsCollection.Find(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	if err = accommodationsCursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodations, nil
}
func (ar *AccommodationRepo) GetByID(id string) (Accommodation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ar.logger.Println(err)
		return Accommodation{}, err
	}

	filter := bson.M{"_id": objID}
	var accommodation Accommodation
	err = accommodationsCollection.FindOne(ctx, filter).Decode(&accommodation)
	if err != nil {
		ar.logger.Println(err)
		return Accommodation{}, err
	}

	return accommodation, nil
}

func (ar *AccommodationRepo) Insert(accommodation *Accommodation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	result, err := accommodationsCollection.InsertOne(ctx, &accommodation)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (ar *AccommodationRepo) Update(id string, accommodation *Accommodation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"name":        accommodation.Name,
		"minGuestNum": accommodation.MinGuestNum,
		"maxGuestNum": accommodation.MaxGuestNum,
		"location":    accommodation.Location,
		"amenities":   accommodation.Amenities,
	}}
	result, err := accommodationsCollection.UpdateOne(ctx, filter, update)
	ar.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	ar.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		ar.logger.Println(err)
		return err
	}
	return nil
}

func (ar *AccommodationRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := accommodationsCollection.DeleteOne(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (ar *AccommodationRepo) SearchAccommodations(searchRequest SearchRequest) (Accommodations, error) {
	// Inicijalizujte context (posle 5 sekundi, prekinite operaciju)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()

	// Pravljenje filtera na osnovu kriterijuma pretrage
	filter := bson.M{
		"location":    searchRequest.Location,
		"minGuestNum": bson.M{"$lte": searchRequest.GuestNum},
		"maxGuestNum": bson.M{"$gte": searchRequest.GuestNum},
	}

	// Izvršavanje upita
	cursor, err := accommodationsCollection.Find(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var accommodations Accommodations
	if err := cursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}

	return accommodations, nil
}

func (ar *AccommodationRepo) getCollection() *mongo.Collection {
	accommodationDatabase := ar.cli.Database("accommodationDB")
	accommodationsCollection := accommodationDatabase.Collection("accommodations")
	return accommodationsCollection
}
func (ar *AccommodationRepo) GetAccommodationsByUserID(userID string) (Accommodations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	// Create a filter to find accommodations owned by the user
	filter := bson.M{"owner._id": objID}

	// Execute the query
	cursor, err := accommodationsCollection.Find(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var accommodations Accommodations
	if err := cursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}

	return accommodations, nil
}

func (ar *AccommodationRepo) DeleteAccommodationsByUserId(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ar.logger.Println(err)
		return err
	}

	filter := bson.M{"owner._id": objID}
	result, err := accommodationsCollection.DeleteMany(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}
