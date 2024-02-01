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

// NoSQL: NotificationRepo struct encapsulating Mongo api client
type NotificationRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*NotificationRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &NotificationRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (nr *NotificationRepo) Disconnect(ctx context.Context) error {
	err := nr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (nr *NotificationRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := nr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		nr.logger.Println(err)
	}

	// Print available databases
	databases, err := nr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		nr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (nr *NotificationRepo) GetAll() (Notifications, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notificationsCollection := nr.getCollection()

	var notifications Notifications
	notificationsCursor, err := notificationsCollection.Find(ctx, bson.M{})
	if err != nil {
		nr.logger.Println(err)
		return nil, err
	}
	if err = notificationsCursor.All(ctx, &notifications); err != nil {
		nr.logger.Println(err)
		return nil, err
	}
	return notifications, nil
}

func (nr *NotificationRepo) GetByID(id string) (Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notificationsCollection := nr.getCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		nr.logger.Println(err)
		return Notification{}, err
	}

	filter := bson.M{"_id": objID}
	var notification Notification
	err = notificationsCollection.FindOne(ctx, filter).Decode(&notification)
	if err != nil {
		nr.logger.Println(err)
		return Notification{}, err
	}

	return notification, nil
}

func (nr *NotificationRepo) Insert(notification *Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notificationsCollection := nr.getCollection()

	result, err := notificationsCollection.InsertOne(ctx, &notification)
	if err != nil {
		nr.logger.Println(err)
		return err
	}
	nr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (nr *NotificationRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notificationsCollection := nr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := notificationsCollection.DeleteOne(ctx, filter)
	if err != nil {
		nr.logger.Println(err)
		return err
	}
	nr.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (nr *NotificationRepo) getCollection() *mongo.Collection {
	notificationDatabase := nr.cli.Database("notificationDB")
	notificationCollection := notificationDatabase.Collection("notifications")
	return notificationCollection
}

func (nr *NotificationRepo) GetNotificationsByUserID(userID string) (Notifications, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notificationCollection := nr.getCollection()
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		nr.logger.Println(err)
		return nil, err
	}
	// Create a filter to find notifications by the user
	filter := bson.M{"host._id": objID}

	// Execute the query
	cursor, err := notificationCollection.Find(ctx, filter)
	if err != nil {
		nr.logger.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications Notifications
	if err := cursor.All(ctx, &notifications); err != nil {
		nr.logger.Println(err)
		return nil, err
	}

	return notifications, nil
}
