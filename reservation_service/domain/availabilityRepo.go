package domain

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// NoSQL: AvailabilityRepo struct encapsulating Mongo api client
type AvailabilityRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*AvailabilityRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &AvailabilityRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (ar *AvailabilityRepo) Disconnect(ctx context.Context) error {
	err := ar.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (ar *AvailabilityRepo) Ping() {
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

func (ar *AvailabilityRepo) GetAll() (AvailabilityPeriods, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	AvailabilityPeriodsCollection := ar.getCollection()

	var availabilityPeriods AvailabilityPeriods
	availabilityPeriodsCursor, err := AvailabilityPeriodsCollection.Find(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	if err = availabilityPeriodsCursor.All(ctx, &availabilityPeriods); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return availabilityPeriods, nil
}

func (ar *AvailabilityRepo) Insert(availabilityPeriod *AvailabilityPeriod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	AvailabilityPeriodsCollection := ar.getCollection()

	result, err := AvailabilityPeriodsCollection.InsertOne(ctx, &availabilityPeriod)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (ar *AvailabilityRepo) Update(id string, availabilityPeriod *AvailabilityPeriod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	AvailabilityPeriodsCollection := ar.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"price":     availabilityPeriod.Price,
		"startDate": availabilityPeriod.StartDate,
		"endDate":   availabilityPeriod.EndDate,
	}}
	result, err := AvailabilityPeriodsCollection.UpdateOne(ctx, filter, update)
	ar.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	ar.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		ar.logger.Println(err)
		return err
	}
	return nil
}

func (ar *AvailabilityRepo) Delete(id string) error {
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

func (ar *AvailabilityRepo) getCollection() *mongo.Collection {
	AvailabilityRepoDatabase := ar.cli.Database("reservationDB")
	AvailabilityRepoCollection := AvailabilityRepoDatabase.Collection("availabilityPeriods")
	return AvailabilityRepoCollection
}
