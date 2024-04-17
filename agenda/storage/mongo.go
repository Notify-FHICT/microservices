package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Notify-FHICT/microservices/agenda/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB defines the interface for interacting with the database
type DB interface {
	//CRUD:

	CreateEvent(event models.Event) error
	ReadEvent(id primitive.ObjectID) (*models.Event, error)
	UpdateEvent(user models.Event) (*models.Event, error)
	DeleteEvent(id primitive.ObjectID) error

	ReadDashboard(id primitive.ObjectID) (*[]models.Event, error)
	LinkNoteID(event models.UpdateNoteID) error
	UnlinkNoteID(event models.UpdateNoteID) error
	LinkTagID(event models.UpdateTagID) error
}

// MongoDB struct implements the DB interface
type MongoDB struct {
	collection *mongo.Collection
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB() (DB, error) {
	// Define MongoDB connection options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:adm1n@noteagenda.5tgrti6.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Select collection
	collection := client.Database("Calendar").Collection("events")

	// Start MongoDB ping to confirm connection
	err = Ping(client)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		collection,
	}, nil
}

// 'Pong' Logs the ping time
func Pong(t time.Time) {
	log.Printf("Pong! (took: %s)", time.Since(t))
}

// 'Ping' Pings MongoDB to confirm connection
func Ping(client *mongo.Client) error {
	defer Pong(time.Now())

	// Sending the actual ping
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	return nil
}

// CRUD operations:
func (db *MongoDB) CreateEvent(event models.Event) error {
	// Insert event into MongoDB collection
	result, err := db.collection.InsertOne(context.TODO(), event)
	fmt.Printf("%s got pushed", result)
	return err
}

func (db MongoDB) ReadEvent(id primitive.ObjectID) (*models.Event, error) {
	// Find event by ID in MongoDB collection
	filter := bson.D{{Key: "_id", Value: id}}
	var result models.Event
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db MongoDB) UpdateEvent(event models.Event) (*models.Event, error) {
	// Update event in MongoDB collection
	filter := bson.D{{Key: "_id", Value: event.ID}}
	var result models.Event
	err := db.collection.FindOneAndReplace(context.TODO(), filter, event).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db MongoDB) DeleteEvent(id primitive.ObjectID) error {
	// Delete event from MongoDB collection
	filter := bson.D{{Key: "_id", Value: id}}
	var result models.Event
	err := db.collection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err
	}
	log.Printf("%s has been deleted", result.ID)
	return err
}

func (db MongoDB) ReadDashboard(id primitive.ObjectID) (*[]models.Event, error) {
	// Read events for a specific user from MongoDB collection
	filter := bson.D{{Key: "userID", Value: id}}
	var result []models.Event
	cursor, err := db.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db MongoDB) LinkNoteID(event models.UpdateNoteID) error {
	// Link note ID to an event in MongoDB collection
	filter := bson.D{{Key: "_id", Value: event.ID}}
	update := bson.D{{"$set", bson.D{{"noteID", event.NoteID}}}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount != 1 || err != nil {
		return err
	}

	return nil
}

func (db MongoDB) UnlinkNoteID(event models.UpdateNoteID) error {
	// Unlink note ID from an event in MongoDB collection
	filter := bson.D{{Key: "noteID", Value: event.NoteID}}
	update := bson.D{{"$set", bson.D{{"noteID", event.ID}}}}

	result, err := db.collection.UpdateMany(context.TODO(), filter, update)

	if result.ModifiedCount != 1 || err != nil {
		return err
	}

	return nil
}

func (db MongoDB) LinkTagID(event models.UpdateTagID) error {
	// Link tag ID to an event in MongoDB collection
	filter := bson.D{{Key: "_id", Value: event.ID}}
	update := bson.D{{"$set", bson.D{{"tagID", event.TagID}}}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount != 1 || err != nil {
		return err
	}

	return nil
}
