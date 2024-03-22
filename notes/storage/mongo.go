package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Notify-FHICT/microservices/notes/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	//CRUD:

	CreateNote(note models.Note) error
	ReadNote(id primitive.ObjectID) (*models.Note, error)
	UpdateNote(note models.Note) (*models.Note, error)
	DeleteNote(id primitive.ObjectID) error
	LinkTagID(event models.UpdateTagID) error
	UpdateContent(note models.Note) error
}

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB() (DB, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:adm1n@noteagenda.5tgrti6.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	collection := client.Database("Notebook").Collection("notes")

	err = Ping(client)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		collection,
	}, nil
}

func Pong(t time.Time) {
	log.Printf("Pong! (took: %s)", time.Since(t))
}

func Ping(client *mongo.Client) error {

	defer Pong(time.Now())

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) CreateNote(note models.Note) error {
	result, err := db.collection.InsertOne(context.TODO(), note)
	fmt.Printf("%s got pushed", result)
	return err
}

func (db MongoDB) ReadNote(id primitive.ObjectID) (*models.Note, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var result models.Note
	err := db.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil

}
func (db MongoDB) UpdateNote(note models.Note) (*models.Note, error) {
	filter := bson.D{{Key: "_id", Value: note.ID}}
	var result models.Note
	err := db.collection.FindOneAndReplace(context.TODO(), filter, note).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db MongoDB) DeleteNote(id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	var result models.Note
	err := db.collection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err
	}
	log.Printf("%s has been deleted", result)
	return err
}

func (db MongoDB) LinkTagID(note models.UpdateTagID) error {
	filter := bson.D{{Key: "_id", Value: note.ID}}
	update := bson.D{{"$set", bson.D{{"tagID", note.TagID}}}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount != 1 || err != nil {
		return err
	}

	return nil
}

func (db MongoDB) UpdateContent(note models.Note) error {
	filter := bson.D{{Key: "_id", Value: note.ID}}
	update := bson.D{{"$set", bson.D{{"content", note.Content}}}}

	result, err := db.collection.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount != 1 || err != nil {
		return err
	}

	return nil
}
